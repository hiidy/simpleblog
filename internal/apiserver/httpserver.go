package apiserver

import (
	"context"

	"github.com/hiidy/simpleblog/internal/pkg/server"
)

type ginServer struct{}

var _ server.Server = (*ginServer)(nil)

func (c *ServerConfig) NewGinServer() server.Server {
	return &ginServer{}
}

func (s *ginServer) RunOrDie() {
	select {}
}

func (s *ginServer) GracefulStop(ctx context.Context) {}
