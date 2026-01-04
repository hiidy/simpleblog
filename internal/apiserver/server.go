package apiserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	handler "github.com/hiidy/simpleblog/internal/apiserver/grpc"
	"github.com/hiidy/simpleblog/internal/pkg/log"
	apiv1 "github.com/hiidy/simpleblog/pkg/api/apiserver/v1"
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
	apiv1.RegisterSimpleBlogServer(grpcServer, handler.NewHandler())
	reflection.Register(grpcServer)

	return &UnionServer{
		cfg: cfg, srv: grpcServer, lis: listener,
	}, nil
}

func (s *UnionServer) Run() error {
	log.Infow("Start to listening the incoming requests on grpc address", "addr", s.cfg.GRPCOptions.Addr)
	go s.srv.Serve(s.lis)

	dialOptions := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(s.cfg.GRPCOptions.Addr, dialOptions...)
	if err != nil {
		return err
	}

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseEnumNumbers: true,
		},
	}))
	if err := apiv1.RegisterSimpleBlogHandler(context.Background(), gwmux, conn); err != nil {
		return err
	}

	log.Infow("Start to listening the incoming requests", "protocol", "http", "addr", s.cfg.HTTPOptions.Addr)
	httpsrv := &http.Server{
		Addr:    s.cfg.HTTPOptions.Addr,
		Handler: gwmux,
	}
	if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
