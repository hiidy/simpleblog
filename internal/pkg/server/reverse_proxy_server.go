package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hiidy/simpleblog/internal/pkg/log"
	genericoptions "github.com/hiidy/simpleblog/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type GRPCGatewayServer struct {
	srv *http.Server
}

func NewGRPCGatewayServer(grpcOptions *genericoptions.GRPCOptions,
	httpOptions *genericoptions.HTTPOptions, registerHandler func(mux *runtime.ServeMux, conn *grpc.ClientConn) error,
) (*GRPCGatewayServer, error) {
	dialOptions := []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 10 * time.Second,
		}),
	}

	dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(grpcOptions.Addr, dialOptions...)
	if err != nil {
		log.Errorw("Failed to dial context", "err", err)
		return nil, err
	}

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseEnumNumbers: true,
		},
	}))
	if err := registerHandler(gwmux, conn); err != nil {
		log.Errorw("Failed to register handler", "err", err)
		return nil, err
	}

	return &GRPCGatewayServer{
		srv: &http.Server{
			Addr:    httpOptions.Addr,
			Handler: gwmux,
		},
	}, nil
}

func (s *GRPCGatewayServer) RunOrDie() {
	log.Infow("Start to listening the incoming requests", "protocol", protocolName(s.srv), "addr", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw("Failed to server HTTP(s) server", "err", err)
	}
}

func (s *GRPCGatewayServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop HTTP(s) server")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Errorw("HTTP(s) server forced to shutdown", "err", err)
	}
}
