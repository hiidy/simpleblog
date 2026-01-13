package apiserver

import (
	"context"
	"net/http"

	handler "github.com/hiidy/simpleblog/internal/apiserver/http"
	mw "github.com/hiidy/simpleblog/internal/pkg/middleware/gin"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/hiidy/simpleblog/internal/pkg/server"
)

type ginServer struct {
	srv server.Server
}

var _ server.Server = (*ginServer)(nil)

func (c *ServerConfig) NewGinServer() server.Server {
	engine := gin.New()
	engine.Use(gin.Recovery(), mw.NoCache, mw.Cors, mw.Secure, mw.RequestIDMiddleware())

	c.InstallRESTAPI(engine)

	httpsrv := server.NewHTTPServer(c.cfg.HTTPOptions, engine)

	return &ginServer{srv: httpsrv}
}

func (c *ServerConfig) InstallRESTAPI(engine *gin.Engine) {
	InstallGenericAPI(engine)

	handler := handler.NewHandler()

	engine.GET("/healthz", handler.Healthz)
}

func InstallGenericAPI(engine *gin.Engine) {
	pprof.Register(engine)

	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, "Page Not Found.")
	})
}

func (s *ginServer) RunOrDie() {
	s.srv.RunOrDie()
}

func (s *ginServer) GracefulStop(ctx context.Context) {
	s.srv.GracefulStop(ctx)
}
