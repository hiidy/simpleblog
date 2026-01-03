package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "github.com/hiidy/simpleblog/pkg/api/apiserver/v1"
)

var (
	addr  = flag.String("addr", "localhost:6666", "The grpc server address to connect to.")
	limit = flag.Int64("limit", 10, "Limit to list users.")
)

func main() {
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to grpc server: %v", err)
	}
	defer conn.Close()

	client := apiv1.NewSimpleBlogClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	resp, err := client.Healthz(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to call healthz: %v", err)
	}

	jsonData, _ := json.Marshal(resp)
	fmt.Println(string(jsonData))
}
