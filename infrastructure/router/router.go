package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
	"github.com/nanmu42/gzip"
)

type router struct {
	Engine  *gin.Engine
	srv     *http.Server
	api     *gin.RouterGroup
	session domain.SessionUsecase
	account domain.AccountUsecase
}

func NewRouter(s domain.SessionUsecase, a domain.AccountUsecase) domain.Router {
	if !config.Env.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.Default()
	e.RemoveExtraSlash = true

	e.Use(gzip.DefaultHandler().Gin)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Env.Port),
		Handler: e,
	}
	return &router{
		srv:     srv,
		Engine:  e,
		api:     e.Group("/api"),
		session: s,
		account: a,
	}
}

func (r *router) Run() {
	go func() {
		Log.Infow("http server running", "port", config.Env.Port)
		if err := r.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Log.Fatalw("http server closed unexpectadly", "err", err)
		}
	}()
}

func (r *router) API() domain.RouteGroup {
	return &ginRouteGroup{r.api, r}
}

func (r *router) Auth(token string) (*domain.Session, *domain.Account, error) {
	session, err := r.session.Fetch(token)
	if err != nil {
		return nil, nil, err
	}
	account, err := r.account.FetchBySession(session)
	if err != nil {
		return nil, nil, err
	}
	return session, account, nil
}

type ginRouteGroup struct {
	*gin.RouterGroup

	router domain.Router
}

func (grg *ginRouteGroup) Group(path string) domain.RouteGroup {
	return &ginRouteGroup{grg.RouterGroup.Group(path), grg.router}
}

func (grg *ginRouteGroup) GET(path string, handler domain.Handler, authorized domain.AuthLevel) {
	grg.RouterGroup.GET(path, grg.endpointHandler(handler, authorized))
}

func (grg *ginRouteGroup) POST(path string, handler domain.Handler, authorized domain.AuthLevel) {
	grg.RouterGroup.POST(path, grg.endpointHandler(handler, authorized))
}

func (grg *ginRouteGroup) PATCH(path string, handler domain.Handler, authorized domain.AuthLevel) {
	grg.RouterGroup.PATCH(path, grg.endpointHandler(handler, authorized))
}

func (grg *ginRouteGroup) DELETE(path string, handler domain.Handler, authorized domain.AuthLevel) {
	grg.RouterGroup.DELETE(path, grg.endpointHandler(handler, authorized))
}

func (grg *ginRouteGroup) endpointHandler(handler domain.Handler, authorized domain.AuthLevel) func(c *gin.Context) {
	if authorized != domain.AuthLevelNone {
		return func(c *gin.Context) {
			context := NewContext(c)

			auth := (context.GetHeader("Authorization"))
			if t, err := context.Cookie("token"); err == nil {
				auth = t
			} else {
				if len(auth) <= 7 || strings.ToLower(auth[:7]) != "bearer " {
					context.Error(domain.ErrUnauthorized.Call())
					return
				}
				auth = auth[7:]
			}
			s, a, err := grg.router.Auth(auth)
			if err != nil {
				context.Error(domain.ErrUnauthorized.Call())
				return
			}
			if a.Type < domain.AccountType(authorized) {
				context.Error(domain.ErrUnauthorized.Call())
				return
			}
			session := NewSession(s, a)

			defer c.Request.Body.Close()
			handler(context, session)
		}
	}
	return func(c *gin.Context) {
		session := EmptySession()
		context := NewContext(c)

		defer c.Request.Body.Close()
		handler(context, session)
	}
}
