package grpc

import (
	apiv1 "github.com/hiidy/simpleblog/pkg/api/apiserver/v1"
)

type Handler struct {
	apiv1.UnimplementedSimpleBlogServer
}

func NewHandler() *Handler {
	return &Handler{}
}
