package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
	"github.com/nanmu42/gzip"
)

type router struct {
	srv    *http.Server
	engine *gin.Engine
	api    *gin.RouterGroup
}

func NewRouter() domain.Router {
	if !config.Env.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.Default()
	e.RemoveExtraSlash = true

	e.Use(gzip.DefaultHandler().Gin)

	// TODO: Add Static with minify
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Env.Port),
		Handler: e,
	}
	return &router{
		srv:    srv,
		engine: e,
		api:    e.Group("/api"),
	}
}

func (r *router) Run() {
	go func() {
		if err := r.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Log.Fatalw("http server closed unexpectadly", "err", err)
		}
	}()
}

func (r *router) API() *gin.RouterGroup {
	return r.api
}

func (r *router) EndpointHandler(handler domain.Handler, authorized bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		session := NewSession(c)
		context := NewContext(c)

		if authorized { // TODO: write authorization
			context.Status(http.StatusUnauthorized)
			return
		}

		defer c.Request.Body.Close()
		handler(context, session)
	}
}
