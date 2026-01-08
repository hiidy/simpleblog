package http

import (
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/hiidy/simpleblog/pkg/api/apiserver/v1"
)

func (h *Handler) Healthz(c *gin.Context) {
	c.JSON(200, &v1.HealthzResponse{
		Status:    v1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	})
}
