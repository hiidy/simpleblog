package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/hiidy/simpleblog/internal/pkg/contextx"
	"github.com/hiidy/simpleblog/internal/pkg/known"
	"github.com/hiidy/simpleblog/pkg/errorsx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		var requestID string
		md, _ := metadata.FromIncomingContext(ctx)

		if reqeustIDs := md[known.XRequestID]; len(reqeustIDs) > 0 {
			requestID = reqeustIDs[0]
		}

		if requestID == "" {
			requestID = uuid.New().String()
			md.Append(known.XRequestID, requestID)
		}

		ctx = metadata.NewIncomingContext(ctx, md)

		_ = grpc.SetHeader(ctx, md)

		ctx = contextx.WithRequestID(ctx, requestID)

		res, err := handler(ctx, req)
		if err != nil {
			return res, errorsx.FromError(err).WithRequestID(requestID)
		}

		return res, nil
	}
}
