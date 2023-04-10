package account_test

import (
	"net/http"
	"testing"

	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/domain/mocks"
	"github.com/minebarteksa/clean-show/test"
	"github.com/minebarteksa/clean-show/usecase/account"
	"github.com/stretchr/testify/assert"
)

func TestPostRegister(t *testing.T) {
	usecase := mocks.NewAccountUsecase(t)
	controller := account.NewAccountController(usecase)

	account := domain.Account{
		DBModel: domain.DBModel{
			ID: 7,
		},
	}
	token := "testingTokenABC123"
	context := &test.MockContext{
		OutCookie: make(map[string]string),
		In:        &domain.AccountCreate{},
	}

	reg := context.In.(*domain.AccountCreate)
	usecase.On("Register", reg).Return(&account, token, nil)

	controller.PostRegister(context, nil)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	assert.Equal(t, token, context.OutCookie["token"])
	data := context.Out.(struct {
		ID    uint   `json:"id"`
		Token string `json:"token"`
	})
	assert.Equal(t, account.ID, data.ID)
	assert.Equal(t, token, data.Token)
}

func TestPostLogin(t *testing.T) {
	usecase := mocks.NewAccountUsecase(t)
	controller := account.NewAccountController(usecase)

	login := domain.AccountLogin{}
	account := domain.Account{DBModel: domain.DBModel{ID: 10}}
	token := "testingTokenABC123"
	context := &test.MockContext{
		OutCookie: make(map[string]string),
		In:        &login,
	}

	usecase.On("Login", &login).Return(&account, token, nil)

	controller.PostLogin(context, nil)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	assert.Equal(t, token, context.OutCookie["token"])
	data := context.Out.(struct {
		ID    uint   `json:"id"`
		Token string `json:"token"`
	})
	assert.Equal(t, account.ID, data.ID)
	assert.Equal(t, token, data.Token)
}

func TestGetByID(t *testing.T) {
	usecase := mocks.NewAccountUsecase(t)
	controller := account.NewAccountController(usecase)

	session := test.NewUserSession(10)
	context := &test.MockContext{
		ParamMap: map[string]string{
			"id": "@me",
		},
	}

	usecase.On("FetchByID", session, uint(10)).Return(session.Account, nil)

	controller.GetByID(context, session)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	data := context.Out.(struct {
		ID uint `json:"id"`
		*domain.Account
	})
	assert.Equal(t, uint(10), data.ID)
	assert.Equal(t, session.Account, data.Account)

	context.ParamMap["id"] = "10"
	context.OutStatus = http.StatusNotImplemented

	controller.GetByID(context, session)

	assert.Equal(t, http.StatusOK, context.OutStatus)
}

func TestPatch(t *testing.T) {
	usecase := mocks.NewAccountUsecase(t)
	controller := account.NewAccountController(usecase)

	account := domain.Account{DBModel: domain.DBModel{ID: 10}}
	session := test.NewUserSession(1)
	context := &test.MockContext{
		ParamMap: map[string]string{
			"id": "10",
		},
		In: map[string]any{},
	}

	usecase.On("Modify", session, uint(10), context.In).Return(&account, nil)

	controller.Patch(context, session)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	data := context.Out.(struct {
		ID uint `json:"id"`
		*domain.Account
	})
	assert.Equal(t, uint(10), data.ID)
	assert.Equal(t, &account, data.Account)
}

func TestGetOrders(t *testing.T) {
	usecase := mocks.NewAccountUsecase(t)
	controller := account.NewAccountController(usecase)

	orders := []domain.Order{{}, {}}
	session := test.NewUserSession(1)
	context := &test.MockContext{
		QueryMap: map[string]string{
			"limit": "5",
			"page":  "2",
		},
		ParamMap: map[string]string{
			"id": "10",
		},
	}

	usecase.On("FetchOrders", session, uint(10), 5, 2).Return(orders, nil)

	controller.GetOrders(context, session)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	data := context.Out.([]domain.Order)
	assert.Equal(t, orders, data)
}

func TestPostPassword(t *testing.T) {
	usecase := mocks.NewAccountUsecase(t)
	controller := account.NewAccountController(usecase)

	password := "newTestPassword123"
	session := test.NewUserSession(1)
	context := &test.MockContext{
		ParamMap: map[string]string{
			"id": "10",
		},
		In: &domain.AccountLogin{
			Password: password,
		},
	}

	usecase.On("ModifyPassword", session, uint(10), password).Return(nil)

	controller.PostPassword(context, session)

	assert.Equal(t, http.StatusNoContent, context.OutStatus)
}

func TestGetLogout(t *testing.T) {
	usecase := mocks.NewAccountUsecase(t)
	controller := account.NewAccountController(usecase)

	session := test.NewUserSession(1)
	context := &test.MockContext{}

	usecase.On("Logout", session).Return(nil)

	controller.GetLogout(context, session)

	assert.Equal(t, http.StatusNoContent, context.OutStatus)
}

func TestDeleteController(t *testing.T) {
	usecase := mocks.NewAccountUsecase(t)
	controller := account.NewAccountController(usecase)

	session := test.NewUserSession(1)
	context := &test.MockContext{
		ParamMap: map[string]string{
			"id": "10",
		},
	}

	usecase.On("Remove", session, uint(10)).Return(nil)

	controller.Delete(context, session)

	assert.Equal(t, http.StatusNoContent, context.OutStatus)
}
