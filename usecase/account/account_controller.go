package account

import (
	"net/http"
	"strconv"

	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/usecase"
)

type accountController struct {
	usecase domain.AccountUsecase
}

func NewAccountController(usecase domain.AccountUsecase) domain.AccountController {
	return &accountController{usecase}
}

func (ac *accountController) Register(router domain.Router) {
	a := router.API().Group("/account")
	a.POST("/register", ac.PostRegister, domain.AuthLevelNone)
	a.POST("/login", ac.PostLogin, domain.AuthLevelNone)
	a.GET("/logout", ac.GetLogout, domain.AuthLevelUser)
	a.GET("/:id", ac.GetByID, domain.AuthLevelUser)
	a.PATCH("/:id", ac.Patch, domain.AuthLevelUser)
	a.GET("/:id/orders", ac.GetOrders, domain.AuthLevelUser)
	a.POST("/:id/password", ac.PostPassword, domain.AuthLevelUser)
	a.DELETE("/:id", ac.Delete, domain.AuthLevelUser)
}

func (ac *accountController) PostRegister(context domain.Context, session domain.UserSession) {
	var register domain.AccountCreate
	err := context.UnmarshalBody(&register)
	if err != nil {
		context.Error(domain.ErrBadRequest.Wrap(err).Call())
		return
	}
	account, token, err := ac.usecase.Register(&register)
	if err != nil {
		context.Error(err)
		return
	}
	context.SetCookie("token", token, 0, "", "", true, true)
	context.JSON(http.StatusOK, struct {
		ID    uint   `json:"id"`
		Token string `json:"token"`
	}{account.ID, token})
}

func (ac *accountController) PostLogin(context domain.Context, session domain.UserSession) {
	var login domain.AccountLogin
	err := context.UnmarshalBody(&login)
	if err != nil {
		context.Error(domain.ErrBadRequest.Wrap(err).Call())
		return
	}
	account, token, err := ac.usecase.Login(&login)
	if err != nil {
		context.Error(err)
		return
	}
	context.SetCookie("token", token, 0, "", "", true, true)
	context.JSON(http.StatusOK, struct {
		ID    uint   `json:"id"`
		Token string `json:"token"`
	}{account.ID, token})
}

func (ac *accountController) GetLogout(context domain.Context, session domain.UserSession) {
	err := ac.usecase.Logout(session)
	if err != nil {
		context.Error(err)
		return
	}
	context.Status(http.StatusNoContent)
}

func (ac *accountController) GetByID(context domain.Context, session domain.UserSession) {
	id := uint(0)
	raw := context.Param("id")
	if raw == "@me" {
		id = session.GetAccountID()
	} else {
		i, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			context.Error(domain.ErrBadRequest.Wrap(err).Call())
			return
		}
		id = uint(i)
	}
	account, err := ac.usecase.FetchByID(session, id)
	if err != nil {
		context.Error(err)
		return
	}
	context.JSON(http.StatusOK,
		struct {
			ID uint `json:"id"`
			*domain.Account
		}{account.ID, account},
	)
}

func (ac *accountController) Patch(context domain.Context, session domain.UserSession) {
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
	account, err := ac.usecase.Modify(session, uint(id), data)
	if err != nil {
		context.Error(err)
		return
	}
	context.JSON(http.StatusOK, struct {
		ID uint `json:"id"`
		*domain.Account
	}{account.ID, account})
}

func (ac *accountController) GetOrders(context domain.Context, session domain.UserSession) {
	limit, page := usecase.GetLimitPage(context)
	rawId := context.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		context.Error(domain.ErrBadRequest.Wrap(err).Call())
		return
	}
	orders, err := ac.usecase.FetchOrders(session, uint(id), limit, page)
	if err != nil {
		context.Error(err)
		return
	}
	context.JSON(http.StatusOK, orders)
}

func (ac *accountController) PostPassword(context domain.Context, session domain.UserSession) {
	rawId := context.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		context.Error(domain.ErrBadRequest.Wrap(err).Call())
		return
	}
	var login domain.AccountLogin
	err = context.UnmarshalBody(&login)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	err = ac.usecase.ModifyPassword(session, uint(id), login.Password)
	if err != nil {
		context.Error(err)
		return
	}
	context.Status(http.StatusNoContent)
}

func (ac *accountController) Delete(context domain.Context, session domain.UserSession) {
	rawId := context.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		context.Error(domain.ErrBadRequest.Wrap(err).Call())
		return
	}
	err = ac.usecase.Remove(session, uint(id))
	if err != nil {
		context.Error(err)
		return
	}
	context.Status(http.StatusNoContent)
}
