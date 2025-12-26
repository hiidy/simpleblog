package apiserver

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"
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
	fmt.Printf("ServerMode from ServerOptions: %s\n", s.cfg.JWTKey)
	fmt.Printf("ServerMode from Viper: %s\n\n", viper.GetString("jwt-key"))

	jsonData, _ := json.MarshalIndent(s.cfg, "", " ")
	fmt.Println(string(jsonData))

	select {}
}
