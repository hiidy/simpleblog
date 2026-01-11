package apiserver

import (
	"context"

	mw "github.com/hiidy/simpleblog/internal/pkg/middleware/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	handler "github.com/hiidy/simpleblog/internal/apiserver/grpc"
	"github.com/hiidy/simpleblog/internal/pkg/server"
	apiv1 "github.com/hiidy/simpleblog/pkg/api/apiserver/v1"
	"google.golang.org/grpc"
)

type grpcServer struct {
	srv  server.Server
	stop func(context.Context)
}

var _ server.Server = (*grpcServer)(nil)

func (c *ServerConfig) NewGRPCServerOr() (server.Server, error) {
	serverOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			mw.RequestIDInterceptor(),
		),
	}

	grpcsrv, err := server.NewGRPCServer(
		c.cfg.GRPCOptions,
		serverOptions,
		func(s grpc.ServiceRegistrar) {
			apiv1.RegisterSimpleBlogServer(s, handler.NewHandler())
		},
	)
	if err != nil {
		return nil, err
	}

	if c.cfg.ServerMode == GRPCServerMode {
		return &grpcServer{
			srv: grpcsrv,
			stop: func(ctx context.Context) {
				grpcsrv.GracefulStop(ctx)
			},
		}, nil
	}

	go grpcsrv.RunOrDie()

	httpsrv, err := server.NewGRPCGatewayServer(
		c.cfg.GRPCOptions,
		c.cfg.HTTPOptions,
		func(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
			return apiv1.RegisterSimpleBlogHandler(context.Background(), mux, conn)
		},
	)
	if err != nil {
		return nil, err
	}

	return &grpcServer{
		srv: httpsrv,
		stop: func(ctx context.Context) {
			grpcsrv.GracefulStop(ctx)
			httpsrv.GracefulStop(ctx)
		},
	}, nil
}

func (s *grpcServer) RunOrDie() {
	s.srv.RunOrDie()
}

func (s *grpcServer) GracefulStop(ctx context.Context) {
	s.stop(ctx)
}
