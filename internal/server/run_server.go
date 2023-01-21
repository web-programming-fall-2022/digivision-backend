package server

import (
	"github.com/arimanius/digivision-backend/internal/bootstrap"
	"github.com/arimanius/digivision-backend/internal/bootstrap/job"
	"github.com/arimanius/digivision-backend/internal/cfg"
	pb "github.com/arimanius/digivision-backend/pkg/api/v1"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RunServer(config cfg.Config) job.WithGracefulShutDown {
	serverRunner, err := bootstrap.NewGrpcServerRunner(config.GrpcServerRunnerConfig)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	grpcServer := serverRunner.GetGrpcServer()
	registerServer(grpcServer)

	go func() {
		if err := serverRunner.Run(); err != nil {
			logrus.Fatal(err.Error())
		}
	}()
	return serverRunner
}

func registerServer(server *grpc.Server) {
	pb.RegisterSearchServiceServer(server, NewSearchServiceServer())
}
