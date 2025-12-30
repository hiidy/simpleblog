package apiserver

import (
	"time"

	"github.com/spf13/viper"

	"github.com/hiidy/simpleblog/internal/pkg/log"
)

type Config struct {
	ServerMode string
	JWTKey     string
	Expiration time.Duration
}

type UnionServer struct {
	cfg *Config
}

func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	return &UnionServer{cfg: cfg}, nil
}

func (s *UnionServer) Run() error {
	log.Infow("ServerMode from ServerOptions", "jwt-key", s.cfg.JWTKey)
	log.Infow("ServerMode from viper", "jwt-key", viper.GetString("jwt-key"))

	select {}
}
