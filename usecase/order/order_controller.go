package order

import (
	"math"
	"net/http"
	"strconv"

	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/usecase"
)

type orderController struct {
	usecase domain.OrderUsecase
}

func NewOrderController(usecase domain.OrderUsecase) domain.OrderController {
	return &orderController{usecase}
}

func (oc *orderController) Register(router domain.Router) {
	o := router.API().Group("/order")
	o.GET("/", oc.Get, domain.AuthLevelStaff)
	o.POST("/", oc.Post, domain.AuthLevelUser)
	o.GET("/:id", oc.GetByID, domain.AuthLevelUser)
	o.PATCH("/:id", oc.Patch, domain.AuthLevelStaff)
	o.POST("/:id/cancel", oc.PostCancel, domain.AuthLevelUser)
	o.DELETE("/:id", oc.Delete, domain.AuthLevelStaff)
}

func (oc *orderController) Get(context domain.Context, session domain.UserSession) {
	limit, page := usecase.GetLimitPage(context)
	count, err := oc.usecase.TotalCount()
	if err != nil {
		context.Error(err)
		return
	}
	list, err := oc.usecase.Fetch(limit, page)
	if err != nil {
		context.Error(err)
		return
	}
	pages := float64(0)
	if limit > 0 {
		pages = math.RoundToEven(float64(count) / float64(limit))
	}
	if pages < 0 {
		pages = 0
	}
	context.JSON(http.StatusOK, domain.DataList[domain.Order]{
		Hits:  count,
		Pages: uint(pages),
		Data:  list,
	})
}

func (oc *orderController) Post(context domain.Context, session domain.UserSession) {
	var create domain.OrderCreate
	err := context.UnmarshalBody(&create)
	if err != nil {
		context.Error(domain.ErrBadRequest.Wrap(err).Call())
		return
	}
	order, err := oc.usecase.Create(session.GetAccountID(), &create)
	if err != nil {
		context.Error(err)
		return
	}
	context.JSON(http.StatusOK,
		struct {
			ID uint `json:"id"`
			*domain.Order
		}{order.ID, order},
	)
}

func (oc *orderController) GetByID(context domain.Context, session domain.UserSession) {
	rawId := context.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		context.Error(domain.ErrBadRequest.Wrap(err).Call())
		return
	}
	order, err := oc.usecase.FetchByID(session, uint(id))
	if err != nil {
		context.Error(err)
		return
	}
	context.JSON(http.StatusOK, order)
}

func (oc *orderController) Patch(context domain.Context, session domain.UserSession) {
	rawId := context.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		context.Error(domain.ErrBadRequest.Wrap(err).Call())
		return
	}
	var data map[string]any
	err = context.UnmarshalBody(&data)
	if err != nil {
		context.Error(domain.ErrBadRequest.Wrap(err).Call())
		return
	}
	err = oc.usecase.Modify(session.GetAccountID(), uint(id), data)
	if err != nil {
		context.Error(err)
		return
	}
	context.Status(http.StatusNoContent)
}

func (oc *orderController) PostCancel(context domain.Context, session domain.UserSession) {
	rawId := context.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		context.Error(domain.ErrBadRequest.Wrap(err).Call())
		return
	}
	err = oc.usecase.Cancel(session, uint(id))
	if err != nil {
		context.Error(err)
		return
	}
	context.Status(http.StatusNoContent)
}

func (oc *orderController) Delete(context domain.Context, session domain.UserSession) {
	rawId := context.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		context.Error(domain.ErrBadRequest.Wrap(err).Call())
		return
	}
	err = oc.usecase.Remove(session.GetAccountID(), uint(id))
	if err != nil {
		context.Error(err)
		return
	}
	context.Status(http.StatusNoContent)
}
