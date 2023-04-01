package product

import (
	"math"
	"net/http"
	"strconv"

	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/usecase"
)

type productController struct {
	usecase domain.ProductUsecase
}

func NewProductController(usecase domain.ProductUsecase) domain.ProductController {
	return &productController{usecase}
}

func (pc *productController) Register(router domain.Router) {
	p := router.API().Group("/product")
	p.GET("/", pc.Get, domain.AuthLevelNone)
	p.POST("/", pc.Post, domain.AuthLevelStaff)
	p.GET("/:id", pc.GetByID, domain.AuthLevelNone)
	p.PATCH("/:id", pc.Patch, domain.AuthLevelStaff)
	p.DELETE("/:id", pc.Delete, domain.AuthLevelStaff)
}

func (pc *productController) Get(context domain.Context, session domain.UserSession) {
	limit, page := usecase.GetLimitPage(context)
	count, err := pc.usecase.TotalCount()
	if err != nil {
		context.Status(http.StatusInternalServerError) // TODO: Better error
		return
	}
	list, err := pc.usecase.Fetch(limit, page)
	if err != nil {
		context.Status(http.StatusInternalServerError) // TODO: Better error
		return
	}
	pages := float64(0)
	if limit > 0 {
		pages = math.RoundToEven(float64(count) / float64(limit))
	}
	if pages < 0 {
		pages = 0
	}
	context.JSON(http.StatusOK, domain.DataList[domain.Product]{
		Hits:  count,
		Pages: uint(pages),
		Data:  list,
	}) // TODO: Add encapsulating struct with product count and page count for the given limit
}

func (pc *productController) Post(context domain.Context, session domain.UserSession) {
	var product domain.Product
	err := context.UnmarshalBody(&product)
	if err != nil {
		context.Status(http.StatusBadRequest) // TODO: Better error
		return
	}
	err = pc.usecase.Create(session.GetAccountID(), &product)
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

func (pc *productController) Patch(context domain.Context, session domain.UserSession) {
	rawId := context.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		context.Status(http.StatusNotFound) // TODO: Better error
		return
	}
	var data map[string]interface{}
	err = context.UnmarshalBody(&data)
	if err != nil {
		context.Status(http.StatusBadRequest) // TODO: Better error
		return
	}
	err = pc.usecase.Modify(session.GetAccountID(), uint(id), data)
	if err != nil {
		context.Status(http.StatusInternalServerError) // TODO: Better error
		return
	}
	context.Status(http.StatusNoContent)
}

func (pc *productController) Delete(context domain.Context, session domain.UserSession) {
	rawId := context.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		context.Status(http.StatusNotFound) // TODO: Better error
		return
	}
	err = pc.usecase.Remove(session.GetAccountID(), uint(id))
	if err != nil {
		context.Status(http.StatusInternalServerError) // TODO: Better error
		return
	}
	context.Status(http.StatusNoContent)
}
