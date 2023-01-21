package bootstrap

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type httpServer struct {
	*http.Server
}

func (server *httpServer) ShutDown(ctx context.Context) error {
	return server.Shutdown(ctx)
}

// newPrometheusServer returns an http server listening on provided address,
// which serves prometheus metrics on /metrics endpoint.
// It also implements job.WithGracefulShutDown.
func newPrometheusServer(address string) *httpServer {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return &httpServer{Server: &http.Server{Addr: address, Handler: mux}}
}
