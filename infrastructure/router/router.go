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
	srv     *http.Server
	engine  *gin.Engine
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

	// TODO: Add Static with minify
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Env.Port),
		Handler: e,
	}
	return &router{
		srv:     srv,
		engine:  e,
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
	return &GinRouteGroup{r.api, r}
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

type GinRouteGroup struct {
	*gin.RouterGroup

	router domain.Router
}

func (grg *GinRouteGroup) Group(path string) domain.RouteGroup {
	return &GinRouteGroup{grg.RouterGroup.Group(path), grg.router}
}

func (grg *GinRouteGroup) GET(path string, handler domain.Handler, authorized domain.AccountType) {
	grg.RouterGroup.GET(path, grg.endpointHandler(handler, authorized))
}

func (grg *GinRouteGroup) POST(path string, handler domain.Handler, authorized domain.AccountType) {
	grg.RouterGroup.POST(path, grg.endpointHandler(handler, authorized))
}

func (grg *GinRouteGroup) PATCH(path string, handler domain.Handler, authorized domain.AccountType) {
	grg.RouterGroup.PATCH(path, grg.endpointHandler(handler, authorized))
}

func (grg *GinRouteGroup) DELETE(path string, handler domain.Handler, authorized domain.AccountType) {
	grg.RouterGroup.DELETE(path, grg.endpointHandler(handler, authorized))
}

func (grg *GinRouteGroup) endpointHandler(handler domain.Handler, authorized domain.AccountType) func(c *gin.Context) {
	if authorized != domain.AccountTypeUnknown {
		return func(c *gin.Context) {
			context := NewContext(c)

			auth := strings.ToLower(context.GetHeader("Authorization"))
			if !strings.HasPrefix(auth, "bearer ") {
				context.String(http.StatusUnauthorized, "401 Unauthorized") // TODO: Better Errors?
				return
			}
			s, a, err := grg.router.Auth(auth[7:])
			if err != nil {
				context.String(http.StatusUnauthorized, "401 Unauthorized") // TODO: Better Errors?
				return
			}
			if authorized == domain.AccountTypeStaff && a.Type != domain.AccountTypeStaff {
				context.String(http.StatusUnauthorized, "401 Unauthorized") // TODO: Better Errors?
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
