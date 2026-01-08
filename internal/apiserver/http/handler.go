package http

import v1 "github.com/hiidy/simpleblog/pkg/api/apiserver/v1"

type Handler struct {
	v1.UnimplementedSimpleBlogServer
}

func NewHandler() *Handler {
	return &Handler{}
}
