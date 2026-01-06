package apiserver

import (
	"time"

	"github.com/hiidy/simpleblog/internal/pkg/log"
	"github.com/hiidy/simpleblog/internal/pkg/server"
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
	HTTPOptions *genericoptions.HTTPOptions
}

type ServerConfig struct {
	cfg *Config
}

type UnionServer struct {
	srv server.Server
}

func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	serverConfig, err := cfg.NewServerConfig()
	if err != nil {
		return nil, err
	}

	log.Infow("Initializing federation server", "server-mode", cfg.ServerMode)

	var srv server.Server
	switch cfg.ServerMode {
	case GinServerMode:
		srv, err = serverConfig.NewGinServer(), nil
	default:
		srv, err = serverConfig.NewGRPCServerOr()
	}

	if err != nil {
		return nil, err
	}
	return &UnionServer{
		srv: srv,
	}, nil
}

func (s *UnionServer) Run() error {
	s.srv.RunOrDie()
	return nil
}

func (cfg *Config) NewServerConfig() (*ServerConfig, error) {
	return &ServerConfig{cfg: cfg}, nil
}
