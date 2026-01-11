package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	v1 "github.com/hiidy/simpleblog/pkg/api/apiserver/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var addr = flag.String("addr", "localhost:6666", "The grpc server address to connect to")

func main() {
	flag.Parse()

	conn, err := grpc.NewClient(*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryClientInterceptor()))
	if err != nil {
		log.Fatalf("Failed to connect to grpc server: %v", err)
	}
	defer conn.Close()

	client := v1.NewSimpleBlogClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	md := metadata.Pairs("custom-header", "value123")
	ctx = metadata.NewOutgoingContext(ctx, md)

	var header metadata.MD
	resp, err := client.Healthz(ctx, nil, grpc.Header(&header))
	if err != nil {
		log.Fatalf("Failed to call healthz: %v", err)
	}

	for key, val := range header {
		fmt.Printf("Response Header (key: %s, value: %s)\n", key, val)
	}

	jsonData, _ := json.Marshal(resp)
	fmt.Println(string(jsonData))
}

func unaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Printf("[UnaryClientInterceptor] Invoking method: %s", method)

		md := metadata.Pairs("interceptor-header", "interceptor-value")
		ctx = metadata.NewOutgoingContext(ctx, md)

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			log.Printf("[UnaryClientInterceptor] Method: %s, Error: %v", method, err)
			return err
		}

		return nil
	}
}
