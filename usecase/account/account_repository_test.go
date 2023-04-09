package account_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/logger"
	"github.com/minebarteksa/clean-show/test"
	"github.com/minebarteksa/clean-show/usecase/account"
	"github.com/stretchr/testify/assert"
)

var (
	repositoryCache domain.AccountRepository
	mockCache       sqlmock.Sqlmock
	preparedCache   []*sqlmock.ExpectedPrepare
)

func NewRepository(t *testing.T) (domain.AccountRepository, sqlmock.Sqlmock, []*sqlmock.ExpectedPrepare) {
	if repositoryCache == nil {
		logger.InitDebug()
		test.SetupConfig()
		db, mock := test.NewMockDB(t)
		mockCache = mock

		preparedCache = []*sqlmock.ExpectedPrepare{
			mock.ExpectPrepare("SELECT \\*\\ FROM accounts WHERE email = \\? AND deleted_at IS NULL"),
			mock.ExpectPrepare("SELECT type, email, name, surname FROM accounts WHERE deleted_at IS NULL AND id = \\?"),
			mock.ExpectPrepare("SELECT \\* FROM accounts WHERE id = \\? AND deleted_at IS NULL"),
			mock.ExpectPrepare("INSERT INTO accounts \\(.*\\) VALUES \\(.*\\) RETURNING id"),
			mock.ExpectPrepare("UPDATE accounts SET .*, updated_at = NOW\\(\\) WHERE id = \\?"),
			mock.ExpectPrepare("UPDATE accounts SET email = '@'\\+email, updated_at = NOW\\(\\) WHERE id = \\?"),
			mock.ExpectPrepare("UPDATE accounts SET deleted_at = NOW\\(\\) WHERE id = \\?"),
		}

		repositoryCache = account.NewAccountRepository(db)
	}
	return repositoryCache, mockCache, preparedCache
}

func TestSelectEMail(t *testing.T) {
	repository, _, prepared := NewRepository(t)
	email := "test@example.com"

	prepared[0].ExpectQuery().WithArgs(email).WillReturnRows(
		test.NewRows("id", "type", "email", "hash", "name", "surname").
			AddRow(7, 2, email, "", "Test", "User"),
	)

	account, err := repository.SelectEMail(email)

	assert.NoError(t, err)
	assert.Equal(t, uint(7), account.ID)
	assert.Equal(t, domain.AccountTypeStaff, account.Type)
	assert.Equal(t, email, account.Email)
}

func TestSelectID(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	prepared[2].ExpectQuery().WithArgs(1).WillReturnRows(
		test.NewRows("id", "type", "email", "name", "surname").
			AddRow(1, 2, "test@example.com", "Test", "User"),
	)

	account, err := repository.SelectID(1, true)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), account.ID)
	assert.Equal(t, domain.AccountTypeStaff, account.Type)
	assert.Empty(t, account.Hash)

	prepared[1].ExpectQuery().WithArgs(1).WillReturnRows(
		test.NewRows("id", "type", "email", "hash", "name", "surname").
			AddRow(1, 2, "test@example.com", "HASH", "Test", "User"),
	)

	account, err = repository.SelectID(1, false)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), account.ID)
	assert.Equal(t, domain.AccountTypeStaff, account.Type)
	assert.Equal(t, "HASH", account.Hash)
}

func TestInsert(t *testing.T) {
	repository, _, prepared := NewRepository(t)
	account := domain.Account{
		Type:    domain.AccountTypeUser,
		Email:   "test@example.com",
		Hash:    "HASH",
		Name:    "Test",
		Surname: "User",
	}

	prepared[3].ExpectQuery().
		WithArgs(account.Type, account.Email, account.Hash, account.Name, account.Surname).
		WillReturnRows(test.IDRow(1))

	err := repository.Insert(&account)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), account.ID)
}

func TestUpdate(t *testing.T) {
	repository, _, prepared := NewRepository(t)
	account := domain.Account{
		DBModel: domain.DBModel{
			ID: 1,
		},
		Type:    domain.AccountTypeUser,
		Email:   "test@example.com",
		Hash:    "HASH",
		Name:    "Test",
		Surname: "User",
	}

	prepared[4].ExpectExec().
		WithArgs(account.Type, account.Email, account.Hash, account.Name, account.Surname, account.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repository.Update(&account)

	assert.NoError(t, err)
}

func TestDeleteRepository(t *testing.T) {
	repository, mock, prepared := NewRepository(t)
	id := uint(7)

	mock.ExpectBegin()
	prepared[5].ExpectExec().WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))
	prepared[6].ExpectExec().WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repository.Delete(id)

	assert.NoError(t, err)
}
