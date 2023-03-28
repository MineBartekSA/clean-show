package product

import (
	"net/http"

	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
)

type productController struct {
	usecase domain.ProductUsecase
}

func NewProductController(usecase domain.ProductUsecase) domain.ProductController {
	return &productController{usecase}
}

func (c *productController) Register(router domain.Router) {
	p := router.API().Group("/product")
	p.GET("/:id", router.EndpointHandler(c.GetByID, false))
	Log.Infow("added product routes")
}

func (c *productController) GetByID(context domain.Context, _ domain.UserSession) {
	//
	context.Status(http.StatusNotImplemented)
}
