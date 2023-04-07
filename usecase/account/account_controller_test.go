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

	context := &test.MockContext{
		OutCookie: make(map[string]string),
		In:        &domain.AccountCreate{},
	}
	account := domain.Account{
		DBModel: domain.DBModel{
			ID: 7,
		},
	}
	token := "testingTokenABC123"

	reg := context.In.(*domain.AccountCreate)
	usecase.On("Register", reg).Return(&account, token, nil)

	controller.PostRegister(context, nil)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	assert.Equal(t, token, context.OutCookie["token"])
	data, ok := context.Out.(struct {
		ID    uint   `json:"id"`
		Token string `json:"token"`
	})
	assert.True(t, ok)
	if ok {
		assert.Equal(t, account.ID, data.ID)
		assert.Equal(t, token, data.Token)
	}
}

// TODO: Write
