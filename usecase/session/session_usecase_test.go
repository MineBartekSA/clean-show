package session_test

import (
	"testing"

	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/domain/mocks"
	"github.com/minebarteksa/clean-show/test"
	"github.com/minebarteksa/clean-show/usecase/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetch(t *testing.T) {
	repository := mocks.NewSessionRepository(t)
	resource, audit := test.NewAuditUsecase(t, domain.ResourceTypeSession)
	usecase := session.NewSessionUsecase(repository, audit)
	token := "testTokenABC"

	repository.On("SelectByToken", token).Return(&domain.Session{
		DBModel: domain.DBModel{
			ID: 1,
		},
		AccountID: 1,
		Token:     token,
	}, nil)
	repository.On("Extend", uint(1)).Return(nil)
	resource.On("Modification", uint(1), uint(1)).Return(nil)

	session, err := usecase.Fetch(token)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), session.AccountID)
	assert.Equal(t, token, session.Token)
}

func TestCreate(t *testing.T) {
	repository := mocks.NewSessionRepository(t)
	resource, audit := test.NewAuditUsecase(t, domain.ResourceTypeSession)
	usecase := session.NewSessionUsecase(repository, audit)

	repository.On("Insert", mock.Anything).Return(nil)
	resource.On("Creation", uint(1), uint(0)).Return(nil)

	session, err := usecase.Create(1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), session.AccountID)
	assert.NotEmpty(t, session.Token)
	assert.Len(t, session.Token, 128)
}

func TestInvalidate(t *testing.T) {
	repository := mocks.NewSessionRepository(t)
	resource, audit := test.NewAuditUsecase(t, domain.ResourceTypeSession)
	usecase := session.NewSessionUsecase(repository, audit)

	repository.On("Delete", uint(1)).Return(nil)
	resource.On("Deletion", uint(1), uint(1)).Return(nil)

	err := usecase.Invalidate(&domain.Session{
		DBModel: domain.DBModel{
			ID: 1,
		},
		AccountID: 1,
	})

	assert.NoError(t, err)
}

func TestInvalidateAccount(t *testing.T) {
	repository := mocks.NewSessionRepository(t)
	resource, audit := test.NewAuditUsecase(t, domain.ResourceTypeSession)
	usecase := session.NewSessionUsecase(repository, audit)

	repository.On("DeleteByAccount", uint(1)).Return(nil)
	resource.On("Deletion", uint(2), uint(0)).Return(nil)

	err := usecase.InvalidateAccount(2, 1)

	assert.NoError(t, err)
}
