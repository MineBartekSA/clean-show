package account

import (
	"net/http"
	"strconv"

	"github.com/minebarteksa/clean-show/domain"
)

type accountController struct {
	usecase domain.AccountUsecase
}

func NewAccountController(usecase domain.AccountUsecase) domain.AccountController {
	return &accountController{usecase}
}

func (ac *accountController) Register(router domain.Router) {
	a := router.API().Group("/account")
	a.GET("/:id", ac.GetByID, domain.AuthLevelUser)
}

func (ac *accountController) GetByID(context domain.Context, session domain.UserSession) {
	id := uint(0)
	raw := context.Param("id")
	if raw == "@me" {
		id = session.GetAccountID()
	} else {
		i, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			context.Status(http.StatusBadRequest) // TODO: Better errors
			return
		}
		id = uint(i)
	}
	account, err := ac.usecase.FetchByID(session, id)
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}
	context.JSON(http.StatusOK,
		struct {
			ID uint `json:"id"`
			*domain.Account
		}{account.ID, account},
	)
}
