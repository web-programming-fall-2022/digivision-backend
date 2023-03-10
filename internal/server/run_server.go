package server

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	img2vecPb "github.com/web-programming-fall-2022/digivision-backend/internal/api/img2vec"
	odPb "github.com/web-programming-fall-2022/digivision-backend/internal/api/od"
	"github.com/web-programming-fall-2022/digivision-backend/internal/bootstrap"
	"github.com/web-programming-fall-2022/digivision-backend/internal/bootstrap/job"
	"github.com/web-programming-fall-2022/digivision-backend/internal/cfg"
	"github.com/web-programming-fall-2022/digivision-backend/internal/img2vec"
	"github.com/web-programming-fall-2022/digivision-backend/internal/od"
	"github.com/web-programming-fall-2022/digivision-backend/internal/productmeta"
	"github.com/web-programming-fall-2022/digivision-backend/internal/rank"
	"github.com/web-programming-fall-2022/digivision-backend/internal/s3"
	"github.com/web-programming-fall-2022/digivision-backend/internal/search"
	"github.com/web-programming-fall-2022/digivision-backend/internal/storage"
	"github.com/web-programming-fall-2022/digivision-backend/internal/token"
	pb "github.com/web-programming-fall-2022/digivision-backend/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math"
	"net/http"
	"regexp"
	"time"
)

func RunServer(ctx context.Context, config cfg.Config) job.WithGracefulShutdown {
	// Create the Img2Vec service
	img2vecConnection, err := grpc.Dial(config.Img2Vec.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			grpcRetry.UnaryClientInterceptor(
				grpcRetry.WithMax(6),
				grpcRetry.WithBackoff(func(attempt uint) time.Duration {
					return 60 * time.Millisecond * time.Duration(math.Pow(3, float64(attempt)))
				}),
				grpcRetry.WithCodes(codes.Unavailable, codes.ResourceExhausted)),
		))
	logrus.Infoln("connection to img2vec service established")
	if err != nil {
		logrus.Fatal(err.Error())
	}
	img2vecClient := img2vecPb.NewImg2VecClient(img2vecConnection)
	i2v := img2vec.NewGrpcImg2Vec(img2vecClient)
	logrus.Infoln("img2vec client created")

	// Create the SearchHandler service
	milvusClient, err := client.NewGrpcClient(ctx, config.Milvus.Addr)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.Infoln("connection to milvus service established")
	searchHandler := search.NewMilvusSearchHandler(
		milvusClient,
		config.Milvus.VectorDim,
		config.Milvus.MetricType,
		config.Milvus.NProbe,
		config.Milvus.CollectionName)
	logrus.Infoln("searchHandler client created")

	// Create the Ranker service
	firstImageRanker := rank.NewFirstImageRanker()
	distCountRanker := rank.NewDistCountRanker()
	rankers := map[pb.Ranker]rank.Ranker{
		pb.Ranker_FIRST_IMAGE: firstImageRanker,
		pb.Ranker_DIST_COUNT:  distCountRanker,
	}
	logrus.Infoln("ranker created")

	// Create the object detector service
	odConnection, err := grpc.Dial(config.ObjectDetector.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			grpcRetry.UnaryClientInterceptor(
				grpcRetry.WithMax(6),
				grpcRetry.WithBackoff(func(attempt uint) time.Duration {
					return 60 * time.Millisecond * time.Duration(math.Pow(3, float64(attempt)))
				}),
				grpcRetry.WithCodes(codes.Unavailable, codes.ResourceExhausted)),
		))
	logrus.Infoln("connection to od service established")
	if err != nil {
		logrus.Fatal(err.Error())
	}
	odClient := odPb.NewObjectDetectorClient(odConnection)
	objectDetector := od.NewGrpcObjectDetector(odClient)
	logrus.Infoln("od client created")

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	httpClient := resty.New()
	fetcher := productmeta.NewDigikalaFetcher(
		"https://www.digikala.com",
		"https://api.digikala.com/v1/product",
		httpClient,
		rdb,
		3,
		5,
	)

	store := storage.NewStorage(&config.MainDB)

	if err := store.Migrate(); err != nil {
		log.Fatal(err)
	}

	tokenManager := token.NewJWTManager(config.JWT.Secret, store, rdb)

	s3Client, err := s3.NewMinioClient(config.S3.Endpoint, config.S3.AccessKey, config.S3.SecretKey, config.S3.UseSSL)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	serverRunner, err := bootstrap.NewGrpcServerRunner(
		config.GrpcServerRunnerConfig,
		[]grpc.UnaryServerInterceptor{NewAuthInterceptor(store, tokenManager).InterceptServer()},
		nil,
	)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	// Create the gRPC server
	grpcServer := serverRunner.GetGrpcServer()

	registerSearchServer(grpcServer, i2v, searchHandler, fetcher, rankers, objectDetector, s3Client, store)

	registerAuthServer(
		grpcServer, tokenManager, store,
		config.JWT.AuthTokenExpire,
		config.JWT.RefreshTokenExpire,
	)

	registerFavoriteServer(
		grpcServer,
		store,
		fetcher,
	)

	go func() {
		logrus.Infoln("Starting grpc server...")
		if err := serverRunner.Run(ctx); err != nil {
			logrus.Fatal(err.Error())
		}
	}()
	return serverRunner
}

func registerSearchServer(
	server *grpc.Server,
	i2v img2vec.Img2Vec,
	searchHandler search.Handler,
	fetcher productmeta.Fetcher,
	rankers map[pb.Ranker]rank.Ranker,
	objectDetector od.ObjectDetector,
	s3Client s3.Client,
	store *storage.Storage,
) {
	pb.RegisterSearchServiceServer(server, NewSearchServiceServer(
		i2v,
		searchHandler,
		fetcher,
		rankers,
		objectDetector,
		s3Client,
		store,
	))
}

func registerAuthServer(
	server *grpc.Server,
	tokenManager token.Manager,
	storage *storage.Storage,
	authTokenExpire int64,
	refreshTokenExpire int64,
) {
	pb.RegisterAuthServiceServer(server, NewAuthServiceServer(
		tokenManager,
		storage,
		authTokenExpire,
		refreshTokenExpire,
	))
}

func registerFavoriteServer(
	server *grpc.Server,
	storage *storage.Storage,
	fetcher productmeta.Fetcher,
) {
	pb.RegisterFavoriteServiceServer(server, NewFavoriteServiceServer(storage, fetcher))
}

func RunHttpServer(ctx context.Context, config cfg.Config) job.WithGracefulShutdown {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			if key == "Authorization" {
				return "x-access-token", true
			}
			return runtime.DefaultHeaderMatcher(key)
		}),
	)
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(32*1024*1024),
			grpc.MaxCallRecvMsgSize(32*1024*1024),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if err := pb.RegisterSearchServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", config.Server.Port), opts); err != nil {
		logrus.Fatal("Failed to start HTTP gateway for search", err.Error())
	}
	if err := pb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", config.Server.Port), opts); err != nil {
		logrus.Fatal("Failed to start HTTP gateway for auth", err.Error())
	}
	if err := pb.RegisterFavoriteServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", config.Server.Port), opts); err != nil {
		logrus.Fatal("Failed to start HTTP gateway for favorite", err.Error())
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpServer.Port),
		Handler: wsproxy.WebsocketProxy(cors(mux)),
	}

	logrus.Info("Starting HTTP/REST Gateway...", srv.Addr)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logrus.Fatal("Failed to start HTTP gateway", err.Error())
		}
	}()
	return srv
}

func allowedOrigin(origin string) bool {
	if viper.GetString("cors") == "*" {
		return true
	}
	if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
		return true
	}
	return false
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowedOrigin(r.Header.Get("Origin")) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
