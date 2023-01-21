package server

import (
	"context"
	"fmt"
	"github.com/arimanius/digivision-backend/internal/bootstrap"
	"github.com/arimanius/digivision-backend/internal/bootstrap/job"
	"github.com/arimanius/digivision-backend/internal/cfg"
	pb "github.com/arimanius/digivision-backend/pkg/api/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func RunServer(ctx context.Context, config cfg.Config) job.WithGracefulShutdown {
	serverRunner, err := bootstrap.NewGrpcServerRunner(config.GrpcServerRunnerConfig)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	grpcServer := serverRunner.GetGrpcServer()
	registerServer(grpcServer)

	go func() {
		if err := serverRunner.Run(ctx); err != nil {
			logrus.Fatal(err.Error())
		}
	}()
	return serverRunner
}

func registerServer(server *grpc.Server) {
	pb.RegisterSearchServiceServer(server, NewSearchServiceServer())
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
