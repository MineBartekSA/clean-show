package router

import (
	"github.com/gin-gonic/gin"
	"github.com/minebarteksa/clean-show/domain"
)

type customContext struct {
	*gin.Context
}

func NewContext(c *gin.Context) domain.Context {
	return &customContext{c}
}
