package apiserver

import (
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	handler "github.com/hiidy/simpleblog/internal/apiserver/grpc"
	"github.com/hiidy/simpleblog/internal/pkg/log"
	v1 "github.com/hiidy/simpleblog/pkg/api/apiserver/v1"
	genericoptions "github.com/hiidy/simpleblog/pkg/options"
)

const (
	GRPCServerMode        = "grpc"
	GRPCGatewayServerMode = "grpc-gateway"
	GinServerMode         = "gin"
)

type Config struct {
	ServerMode  string
	JWTKey      string
	Expiration  time.Duration
	GRPCOptions *genericoptions.GRPCOptions
}

type UnionServer struct {
	cfg *Config
	srv *grpc.Server
	lis net.Listener
}

func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	listener, err := net.Listen("tcp", cfg.GRPCOptions.Addr)
	if err != nil {
		log.Fatalw("Failed to listen", "err", err)
		return nil, err
	}

	grpcServer := grpc.NewServer()
	v1.RegisterSimpleBlogServer(grpcServer, handler.NewHandler())
	reflection.Register(grpcServer)

	return &UnionServer{
		cfg: cfg, srv: grpcServer, lis: listener,
	}, nil
}

func (s *UnionServer) Run() error {
	log.Infow("Start to listening the incoming requests on grpc address", "addr", s.cfg.GRPCOptions.Addr)
	return s.srv.Serve(s.lis)
}
