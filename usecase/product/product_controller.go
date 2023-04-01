package product

import (
	"net/http"
	"strconv"

	"github.com/minebarteksa/clean-show/domain"
)

type productController struct {
	usecase domain.ProductUsecase
}

func NewProductController(usecase domain.ProductUsecase) domain.ProductController {
	return &productController{usecase}
}

func (pc *productController) Register(router domain.Router) {
	p := router.API().Group("/product")
	p.POST("/", pc.CreateNew, domain.AuthLevelStaff)
	p.GET("/:id", pc.GetByID, domain.AuthLevelNone)
}

func (pc *productController) CreateNew(context domain.Context, session domain.UserSession) {
	var product domain.Product
	err := context.UnmarshalBody(&product)
	if err != nil {
		context.Status(http.StatusBadRequest) // TODO: Better error
		return
	}
	err = pc.usecase.Create(session.GetAccount().ID, &product)
	if err != nil {
		context.Status(http.StatusInternalServerError) // TODO: Better error
		return
	}
	context.JSON(http.StatusOK,
		struct {
			ID uint `json:"id"`
			*domain.Product
		}{product.ID, &product},
	)
}

func (pc *productController) GetByID(context domain.Context, _ domain.UserSession) {
	rawId := context.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		context.Status(http.StatusNotFound) // TODO: Better error
		return
	}
	product, err := pc.usecase.FetchByID(uint(id))
	if err != nil {
		context.Status(http.StatusNotFound) // TODO: Better error
		return
	}
	context.JSON(http.StatusOK, product)
}
