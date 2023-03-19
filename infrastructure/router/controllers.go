package router

import (
	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
)

type appController struct {
	controllers []domain.Controller
}

func NewAppController(controllers []domain.Controller) domain.Controller {
	return &appController{controllers}
}

func (c *appController) Register(router domain.Router) {
	Log.Infow("registering controllers", "len", len(c.controllers))
	for _, controller := range c.controllers {
		controller.Register(router)
	}
}
