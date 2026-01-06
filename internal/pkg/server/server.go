package server

import (
	"context"
	"net/http"
)

type Server interface {
	RunOrDie()

	GracefulStop(ctx context.Context)
}

func protocolName(server *http.Server) string {
	if server.TLSConfig != nil {
		return "https"
	}
	return "http"
}
