package domain

import (
	"github.com/gin-gonic/gin"
)

type Handler func(context Context, session UserSession)

type Router interface {
	Run()
	API() *gin.RouterGroup
	EndpointHandler(handler Handler, authorized bool) func(*gin.Context)
}

type Controller interface {
	Register(router Router)
}
