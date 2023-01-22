package server

import (
	"context"
	"fmt"
	img2vecPb "github.com/arimanius/digivision-backend/internal/api/img2vec"
	odPb "github.com/arimanius/digivision-backend/internal/api/od"
	"github.com/arimanius/digivision-backend/internal/bootstrap"
	"github.com/arimanius/digivision-backend/internal/bootstrap/job"
	"github.com/arimanius/digivision-backend/internal/cfg"
	"github.com/arimanius/digivision-backend/internal/img2vec"
	"github.com/arimanius/digivision-backend/internal/od"
	"github.com/arimanius/digivision-backend/internal/rank"
	"github.com/arimanius/digivision-backend/internal/search"
	pb "github.com/arimanius/digivision-backend/pkg/api/v1"
	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"math"
	"net/http"
	"time"
)

func RunServer(ctx context.Context, config cfg.Config) job.WithGracefulShutdown {
	serverRunner, err := bootstrap.NewGrpcServerRunner(config.GrpcServerRunnerConfig)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	// Create the gRPC server
	grpcServer := serverRunner.GetGrpcServer()

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
	if err != nil {
		logrus.Fatal(err.Error())
	}
	defer func(img2vecConnection *grpc.ClientConn) {
		err := img2vecConnection.Close()
		if err != nil {
			logrus.Fatal(err.Error())
		}
	}(img2vecConnection)
	img2vecClient := img2vecPb.NewImg2VecClient(img2vecConnection)
	i2v := img2vec.NewGrpcImg2Vec(img2vecClient)

	// Create the SearchHandler service
	milvusClient, err := client.NewGrpcClient(ctx, config.Milvus.Addr)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	searchHandler := search.NewMilvusSearchHandler(
		milvusClient,
		config.Milvus.VectorDim,
		config.Milvus.MetricType,
		config.Milvus.NProbe,
		config.Milvus.CollectionName)

	// Create the Ranker service
	ranker := rank.NewFirstImageRanker()

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
	if err != nil {
		logrus.Fatal(err.Error())
	}
	defer func(odConnection *grpc.ClientConn) {
		err := odConnection.Close()
		if err != nil {
			logrus.Fatal(err.Error())
		}
	}(odConnection)
	odClient := odPb.NewObjectDetectorClient(odConnection)
	objectDetector := od.NewGrpcObjectDetector(odClient)

	registerServer(grpcServer, i2v, searchHandler, ranker, objectDetector)

	go func() {
		if err := serverRunner.Run(ctx); err != nil {
			logrus.Fatal(err.Error())
		}
	}()
	return serverRunner
}

func registerServer(server *grpc.Server, i2v img2vec.Img2Vec, searchHandler search.Handler, ranker rank.Ranker, objectDetector od.ObjectDetector) {
	pb.RegisterSearchServiceServer(server, NewSearchServiceServer(i2v, searchHandler, ranker, objectDetector))
}

func RunHttpServer(ctx context.Context, config cfg.Config) job.WithGracefulShutdown {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterSearchServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", config.Server.Port), opts); err != nil {
		logrus.Fatal("Failed to start HTTP gateway", err.Error())
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpServer.Port),
		Handler: mux,
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
