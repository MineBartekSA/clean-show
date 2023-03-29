package product

import (
	"log"
	"net/http"

	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
)

type productController struct {
	usecase domain.ProductUsecase
	audit   domain.AuditUsecase
}

func NewProductController(usecase domain.ProductUsecase, audit domain.AuditUsecase) domain.ProductController {
	return &productController{usecase, audit}
}

func (c *productController) Register(router domain.Router) {
	p := router.API().Group("/product")
	p.GET("/:id", c.GetByID, false)
	Log.Infow("added product routes")
}

func (c *productController) GetByID(context domain.Context, _ domain.UserSession) {
	err := c.audit.Creation(0, domain.ResourceTypeOrder, 0)
	if err != nil {
		log.Panicln(err)
	}
	context.Status(http.StatusOK)
}
