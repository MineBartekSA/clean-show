package account_test

import (
	"testing"

	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/domain/mocks"
	"github.com/minebarteksa/clean-show/test"
	"github.com/minebarteksa/clean-show/usecase/account"
	"github.com/stretchr/testify/assert"
)

func AccountAudit(t *testing.T) (*mocks.AuditResource, *mocks.AuditResource, domain.AuditUsecase) {
	mock := mocks.NewAuditUsecase(t)
	account := mocks.NewAuditResource(t)
	password := mocks.NewAuditResource(t)
	mock.On("Resource", domain.ResourceTypeAccount).Return(account)
	mock.On("Resource", domain.ResourceTypeAccountPassword).Return(password)
	return account, password, mock
}

func TestLogin(t *testing.T) {
	repository := mocks.NewAccountRepository(t)
	session := mocks.NewSessionUsecase(t)
	hasher := test.NewMockHasher().Add("test", "HASH")
	_, _, audit := AccountAudit(t)
	usecase := account.NewAccountUsecase(repository, nil, session, audit, hasher)

	login := domain.AccountLogin{
		Email:    "test@example.com",
		Password: "test",
	}
	account := domain.Account{
		DBModel: domain.DBModel{
			ID: 7,
		},
		Type:    domain.AccountTypeUser,
		Email:   login.Email,
		Hash:    "HASH",
		Name:    "Test",
		Surname: "User",
	}
	token := "testingToken123ABC"

	repository.On("SelectEMail", login.Email).Return(&account, nil)
	session.On("Create", account.ID).Return(&domain.Session{
		AccountID: account.ID,
		Token:     token,
	}, nil)

	laccount, ltoken, err := usecase.Login(&login)

	assert.NoError(t, err)
	assert.Equal(t, account, *laccount)
	assert.Equal(t, token, ltoken)
}

func TestRegister(t *testing.T) {
	repository := mocks.NewAccountRepository(t)
	session := mocks.NewSessionUsecase(t)
	hasher := test.NewMockHasher().Add("test", "HASH")
	_, _, audit := AccountAudit(t)
	usecase := account.NewAccountUsecase(repository, nil, session, audit, hasher)

	register := domain.AccountCreate{
		AccountLogin: &domain.AccountLogin{
			Email:    "test@example.com",
			Password: "test",
		},
		Name:    "Test",
		Surname: "User",
	}
	account := domain.Account{
		Type:    domain.AccountTypeUser,
		Email:   register.Email,
		Hash:    "HASH",
		Name:    register.Name,
		Surname: register.Surname,
	}
	token := "testingToken123ABC"

	repository.On("Insert", &account).Return(nil)
	session.On("Create", account.ID).Return(&domain.Session{
		AccountID: account.ID,
		Token:     token,
	}, nil)

	raccount, rtoken, err := usecase.Register(&register)

	assert.NoError(t, err)
	assert.Equal(t, account, *raccount)
	assert.Equal(t, token, rtoken)
}

func TestFetchBySession(t *testing.T) {
	repository := mocks.NewAccountRepository(t)
	_, _, audit := AccountAudit(t)
	usecase := account.NewAccountUsecase(repository, nil, nil, audit, nil)

	account := domain.Account{
		DBModel: domain.DBModel{
			ID: 7,
		},
	}
	session := domain.Session{
		AccountID: account.ID,
		Token:     "test",
	}

	repository.On("SelectID", session.AccountID, true).Return(&account, nil)

	faccount, err := usecase.FetchBySession(&session)

	assert.NoError(t, err)
	assert.Equal(t, account, *faccount)
}

func TestFetchByID(t *testing.T) {
	repository := mocks.NewAccountRepository(t)
	_, _, audit := AccountAudit(t)
	usecase := account.NewAccountUsecase(repository, nil, nil, audit, nil)

	account := domain.Account{
		DBModel: domain.DBModel{
			ID: 7,
		},
		Type:    domain.AccountTypeStaff,
		Email:   "test@example.com",
		Name:    "Test",
		Surname: "User",
	}
	session := test.MockUserSession{
		Account: &account,
	}

	repository.On("SelectID", account.ID, false).Return(&account, nil)

	faccount, err := usecase.FetchByID(&session, 7)

	assert.NoError(t, err)
	assert.Equal(t, account, *faccount)
}

func TestModify(t *testing.T) {
	repository := mocks.NewAccountRepository(t)
	auditAccount, _, audit := AccountAudit(t)
	usecase := account.NewAccountUsecase(repository, nil, nil, audit, nil)

	account := domain.Account{
		DBModel: domain.DBModel{
			ID: 7,
		},
	}
	session := test.MockUserSession{
		Account: &account,
	}

	repository.On("SelectID", account.ID, false).Return(&account, nil)
	repository.On("Update", &account).Return(nil)
	auditAccount.On("Modification", session.GetAccountID(), account.ID).Return(nil)

	err := usecase.Modify(&session, 7, map[string]any{
		"email": "test@test.com",
		"name":  "A",
		"hash":  "hh",
	})

	assert.NoError(t, err)
	assert.Equal(t, "test@test.com", account.Email)
	assert.Equal(t, "A", account.Name)
	assert.Empty(t, account.Hash)
}

func TestFetchOrders(t *testing.T) {
	repository := mocks.NewAccountRepository(t)
	order := mocks.NewOrderUsecase(t)
	_, _, audit := AccountAudit(t)
	usecase := account.NewAccountUsecase(repository, order, nil, audit, nil)

	session := test.MockUserSession{
		Account: &domain.Account{
			DBModel: domain.DBModel{
				ID: 7,
			},
		},
	}

	order.On("FetchByAccount", session.Account.ID).Return([]domain.Order{{}, {}, {}}, nil)

	orders, err := usecase.FetchOrders(&session, uint(7))

	assert.NoError(t, err)
	assert.Len(t, orders, 3)
}

func TestModifyPassword(t *testing.T) {
	repository := mocks.NewAccountRepository(t)
	hasher := test.NewMockHasher().Add("hello", "HASH")
	_, password, audit := AccountAudit(t)
	usecase := account.NewAccountUsecase(repository, nil, nil, audit, hasher)

	session := test.MockUserSession{
		Account: &domain.Account{
			DBModel: domain.DBModel{
				ID: 7,
			},
		},
	}

	repository.On("SelectID", session.Account.ID, false).Return(session.Account, nil)
	repository.On("Update", session.Account).Return(nil)
	password.On("Modification", session.Account.ID, session.Account.ID).Return(nil)

	err := usecase.ModifyPassword(&session, 7, "hello")

	assert.NoError(t, err)
	assert.Equal(t, "HASH", session.Account.Hash)
}

func TestLogout(t *testing.T) {
	repository := mocks.NewAccountRepository(t)
	session := mocks.NewSessionUsecase(t)
	_, _, audit := AccountAudit(t)
	usecase := account.NewAccountUsecase(repository, nil, session, audit, nil)

	userSession := test.MockUserSession{
		Session: &domain.Session{},
	}

	session.On("Invalidate", userSession.Session).Return(nil)

	err := usecase.Logout(&userSession)

	assert.NoError(t, err)
}

func TestRemove(t *testing.T) {
	repository := mocks.NewAccountRepository(t)
	session := mocks.NewSessionUsecase(t)
	order := mocks.NewOrderUsecase(t)
	auditAccount, _, audit := AccountAudit(t)
	usecase := account.NewAccountUsecase(repository, order, session, audit, nil)

	userSession := test.MockUserSession{
		Account: &domain.Account{
			DBModel: domain.DBModel{
				ID: 7,
			},
		},
	}

	repository.On("Delete", userSession.Account.ID).Return(nil)
	order.On("CancelByAccount", userSession.Account.ID, userSession.Account.ID).Return(nil)
	session.On("InvalidateAccount", userSession.Account.ID, userSession.Account.ID).Return(nil)
	auditAccount.On("Deletion", userSession.Account.ID, userSession.Account.ID).Return(nil)

	err := usecase.Remove(&userSession, userSession.Account.ID)

	assert.NoError(t, err)
}
