package router_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/domain/mocks"
	"github.com/minebarteksa/clean-show/infrastructure/router"
	"github.com/minebarteksa/clean-show/logger"
	"github.com/minebarteksa/clean-show/test"
	"github.com/stretchr/testify/assert"
)

func TestRouterAuthorization(t *testing.T) {
	test.SetupConfig()
	logger.InitDebug()
	session := mocks.NewSessionUsecase(t)
	account := mocks.NewAccountUsecase(t)
	r := router.NewRouter(session, account)
	engine := reflect.ValueOf(r).Elem().FieldByName("Engine").Interface().(*gin.Engine)

	rec := httptest.NewRecorder()
	r.API().GET("/test", func(context domain.Context, session domain.UserSession) {
		assert.True(t, session.Authorized())
		assert.True(t, !session.IsStaff())
		context.Status(http.StatusOK)
	}, domain.AuthLevelUser)
	req := httptest.NewRequest("GET", "/api/test", nil)

	// No Authorization header

	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Result().StatusCode)

	// Invalid authorization token

	session.On("Fetch", "testingToken123").Return(nil, domain.ErrUnauthorized.Call())

	rec = httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer testingToken123")
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Result().StatusCode)

	// Ok

	s := domain.Session{}
	session.On("Fetch", "testingToken123@").Return(&s, nil)
	a := domain.Account{
		Type: domain.AccountTypeUser,
	}
	account.On("FetchBySession", &s).Return(&a, nil)

	rec = httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer testingToken123@")
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	// Not staff

	r.API().GET("/test_staff", func(context domain.Context, session domain.UserSession) {
		assert.True(t, session.Authorized())
		context.Status(http.StatusOK)
	}, domain.AuthLevelStaff)
	req = httptest.NewRequest("GET", "/api/test_staff", nil)
	req.Header.Set("Authorization", "Bearer testingToken123@")

	rec = httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Result().StatusCode)
}
