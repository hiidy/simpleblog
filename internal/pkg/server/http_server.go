package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/hiidy/simpleblog/internal/pkg/log"
	genericoptions "github.com/hiidy/simpleblog/pkg/options"
)

type HTTPServer struct {
	srv *http.Server
}

func NewHTTPServer(httpOptions *genericoptions.HTTPOptions, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		srv: &http.Server{
			Addr:    httpOptions.Addr,
			Handler: handler,
		},
	}
}

func (s *HTTPServer) RunOrDie() {
	log.Infow("Start to listening the incoming requests", "protocol", protocolName(s.srv), "addr", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw("HTTP server failed", "err", err)
	}
}

func (s *HTTPServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop HTTP server")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Errorw("HTTP server forced to shutdown", "err", err)
	}
}
