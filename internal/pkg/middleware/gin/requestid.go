package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hiidy/simpleblog/internal/pkg/contextx"
	"github.com/hiidy/simpleblog/internal/pkg/known"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.Request.Header.Get(known.XRequestID)

		if requestID == "" {
			requestID = uuid.NewString()
		}

		newCtx := contextx.WithRequestID(ctx.Request.Context(), requestID)
		ctx.Request = ctx.Request.WithContext(newCtx)

		ctx.Writer.Header().Set(known.XRequestID, requestID)

		ctx.Next()
	}
}
