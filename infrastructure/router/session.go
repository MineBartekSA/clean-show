package router

import (
	"github.com/gin-gonic/gin"
	"github.com/minebarteksa/clean-show/domain"
)

type session struct {
	//
}

func NewSession(c *gin.Context) domain.Session {
	return &session{}
}
