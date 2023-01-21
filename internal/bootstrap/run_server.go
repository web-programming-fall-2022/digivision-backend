package bootstrap

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/arimanius/digivision-backend/internal/bootstrap/job"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// GrpcServerRunner is used to painlessly run a gRPC sever. Use GetGrpcServer to register your service(s).
// Shutdown will shut it down gracefully.
type GrpcServerRunner interface {
	GetGrpcServer() *grpc.Server
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

func NewGrpcServerRunner(config GrpcServerRunnerConfig) (GrpcServerRunner, error) {
	runner := &grpcServerRunner{
		config:           &config,
		shutDownReqChan:  make(chan bool),
		shutDownDoneChan: make(chan bool),
	}
	err := runner.initialize()
	if err != nil {
		return nil, err
	}
	return runner, nil
}

// GrpcServerRunnerConfig is used to configure GrpcServerRunnerConfig behaviour. It can be used with viper.
// You can provide your own net.Listener implementation instead of host and port, by setting Server.Port to 0 and
// setting the value of Server.Connection.
type GrpcServerRunnerConfig struct {
	Server struct {
		Connection net.Listener `mapstructure:",omitempty"`
		Host       string
		Port       int
		Auth       struct {
			CertFile string `mapstructure:"cert-file" yaml:"cert-file"`
			KeyFile  string `mapstructure:"key-file" yaml:"key-file"`
		}
	}
	Prometheus struct {
		Enabled bool
		Host    string
		Port    int
		Prefix  string
		Buckets []float64
	}
}

// internal and implementation

func (config *GrpcServerRunnerConfig) grpcAddress() string {
	if config.Server.Port == 0 {
		return fmt.Sprint(config.Server.Connection)
	}
	return fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
}

type grpcServerRunner struct {
	config *GrpcServerRunnerConfig

	shutDownReqChan  chan bool
	shutDownDoneChan chan bool
	shutDownError    error

	netListener net.Listener
	grpcServer  *grpc.Server
}

type terminableResources struct {
	netListener net.Listener
	grpcServer  *grpc.Server
	jobs        []job.WithGracefulShutdown
}

func (r *grpcServerRunner) GetGrpcServer() *grpc.Server {
	return r.grpcServer
}

func (r *grpcServerRunner) Shutdown(ctx context.Context) error {
	r.shutDownReqChan <- true
	select {
	case <-r.shutDownDoneChan:
		return r.shutDownError
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *grpcServerRunner) initialize() error {
	if r.config.Server.Port != 0 {
		conn, err := net.Listen("tcp", r.config.grpcAddress())
		if err != nil {
			logrus.WithError(err).Errorf("Failed to listen on %q", r.config.grpcAddress())
			return err
		}
		r.netListener = conn
	} else {
		if r.config.Server.Connection == nil {
			return errors.New("either Server.Port should be not-zero or Server.Connection should be non-empty")
		}
		r.netListener = r.config.Server.Connection
	}

	opts, err := r.serverOptions()
	if err != nil {
		r.closeConnection()
		return err
	}

	grpcServer := grpc.NewServer(opts...)
	healthCheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthCheck)
	r.grpcServer = grpcServer
	return nil
}

func (r *grpcServerRunner) closeConnection() {
	if r.config.Server.Port != 0 {
		if err := r.netListener.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close connection")
		}
	}
}

func (r *grpcServerRunner) Run(ctx context.Context) (runErr error) {
	var terminableJobs []job.WithGracefulShutdown

	go func() {
		logrus.Infof("GRPC Listening on %s", r.config.grpcAddress())
		err := r.grpcServer.Serve(r.netListener)
		if err != nil {
			logrus.WithError(err).Error("failed to grpcServer.Serve(conn)")
			runErr = err
		}
	}()
	if runErr != nil {
		r.closeConnection()
		return
	}

	if r.config.Prometheus.Enabled {
		grpc_prometheus.EnableHandlingTimeHistogram(grpc_prometheus.WithHistogramBuckets(r.config.Prometheus.Buckets))
		promAddress := fmt.Sprintf("%v:%v", r.config.Prometheus.Host, r.config.Prometheus.Port)
		promServer := newPrometheusServer(promAddress)

		terminableJobs = append(terminableJobs, promServer)

		go func() {
			logrus.Infof("Prometheus Listening on %s", promAddress)
			err := promServer.ListenAndServe()
			if err != nil {
				logrus.WithError(err).Errorf("failed to serve prometheus on %v", promAddress)
				runErr = err
			}
		}()
		if runErr != nil {
			r.closeConnection()
			return
		}
	}

	terminableResources := terminableResources{
		grpcServer: r.grpcServer,
		jobs:       terminableJobs,
	}
	if r.config.Server.Port != 0 {
		terminableResources.netListener = r.netListener
	}
	r.waitForTermination(ctx, terminableResources)
	return
}

func (r *grpcServerRunner) serverOptions() ([]grpc.ServerOption, error) {
	opts := getGrpcMiddlewares(r.config)
	certFile := r.config.Server.Auth.CertFile
	keyFile := r.config.Server.Auth.KeyFile
	if certFile == "" || keyFile == "" {
		logrus.Info("No credential provided!, running in insecure mode.")
		return opts, nil
	}
	creds, err := r.newServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return append(opts, grpc.Creds(creds)), nil
}

func getGrpcMiddlewares(config *GrpcServerRunnerConfig) []grpc.ServerOption {
	panicRecoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			var l *logrus.Entry
			if err, ok := p.(error); ok {
				l = logrus.WithError(err)
			} else {
				l = logrus.WithField("p", err)
			}
			l.Errorln("Panic when handling request")
			return status.Error(codes.Unknown, "panic triggered")
		}),
	}
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	unaryInterceptors = append(unaryInterceptors, grpc_recovery.UnaryServerInterceptor(panicRecoveryOpts...))
	streamInterceptors = append(streamInterceptors, grpc_recovery.StreamServerInterceptor(panicRecoveryOpts...))

	if config.Prometheus.Enabled {
		unaryInterceptors = append(unaryInterceptors, grpc_prometheus.UnaryServerInterceptor)
		streamInterceptors = append(streamInterceptors, grpc_prometheus.StreamServerInterceptor)
	}

	return []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(unaryInterceptors...),
		grpc_middleware.WithStreamServerChain(streamInterceptors...),
	}
}

func (r *grpcServerRunner) newServerTLSFromFile(certFile, keyFile string) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		logrus.WithError(err).Error("Failed to create credentials")
		return nil, err
	}
	return credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}), nil
}

func (r *grpcServerRunner) shutdownServices(ctx context.Context, resources terminableResources) {
	logrus.Info("Stopping GRPC server and all jobs.")
	r.runAndWait(
		func() {
			resources.grpcServer.GracefulStop()
			if resources.netListener != nil {
				if err := resources.netListener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
					logrus.WithField("operation", "close net listener").Error(err.Error())
					r.shutDownError = err
				}
			}
		},
		func() {
			if err := job.Shutdown(ctx, resources.jobs, 15*time.Second); err != nil {
				if r.shutDownError != nil {
					logrus.Error(err.Error())
				} else {
					r.shutDownError = err
				}
			}
		},
	)
	r.shutDownDoneChan <- true
}

func (r *grpcServerRunner) runAndWait(functions ...func()) {
	wg := &sync.WaitGroup{}
	for i := range functions {
		wg.Add(1)
		fn := functions[i] // Beware of https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		go func() {
			fn()
			wg.Done()
		}()
	}
	wg.Wait()
}

func (r *grpcServerRunner) waitForTermination(ctx context.Context, resources terminableResources) {
	<-r.shutDownReqChan
	r.shutdownServices(ctx, resources)
}
