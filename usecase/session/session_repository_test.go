package session_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/logger"
	"github.com/minebarteksa/clean-show/test"
	"github.com/minebarteksa/clean-show/usecase/session"
	"github.com/stretchr/testify/assert"
)

var (
	repositoryCache domain.SessionRepository
	mockCache       sqlmock.Sqlmock
	preparedCache   []*sqlmock.ExpectedPrepare
)

func NewRepository(t *testing.T) (domain.SessionRepository, sqlmock.Sqlmock, []*sqlmock.ExpectedPrepare) {
	if repositoryCache == nil {
		logger.InitDebug()
		test.SetupConfig()
		db, mock := test.NewMockDB(t)
		mockCache = mock

		preparedCache = []*sqlmock.ExpectedPrepare{
			mock.ExpectPrepare("SELECT \\* FROM sessions WHERE updated_at > NOW\\(\\) \\+ INTERVAL '-30 MINUTE' AND token = \\? AND deleted_at IS NULL"),
			mock.ExpectPrepare("INSERT INTO sessions \\(.*\\) VALUES \\(.*\\) RETURNING id"),
			mock.ExpectPrepare("UPDATE sessions SET updated_at = NOW\\(\\) WHERE id = \\?"),
			mock.ExpectPrepare("UPDATE sessions SET deleted_at = NOW\\(\\) WHERE id = \\?"),
			mock.ExpectPrepare("UPDATE sessions SET deleted_at = NOW\\(\\) WHERE account_id = \\?"),
		}

		repositoryCache = session.NewSessionRepository(db)
	}
	return repositoryCache, mockCache, preparedCache
}

func TestSelectByToken(t *testing.T) {
	repository, _, prepared := NewRepository(t)
	token := "testTokenABC"

	prepared[0].ExpectQuery().WithArgs(token).WillReturnRows(test.NewRows("id", "account_id", "token").AddRow(1, 1, token))

	session, err := repository.SelectByToken(token)

	assert.NoError(t, err)
	assert.Equal(t, token, session.Token)
}

func TestInsert(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	session := domain.Session{
		AccountID: 1,
		Token:     "testTokenABC",
	}

	prepared[1].ExpectQuery().WithArgs(session.AccountID, session.Token).WillReturnRows(test.IDRow(1))

	err := repository.Insert(&session)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), session.ID)
}

func TestExtend(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	prepared[2].ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repository.Extend(1)

	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	prepared[3].ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repository.Delete(1)

	assert.NoError(t, err)
}

func TestDeleteByAccount(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	prepared[4].ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 5))

	err := repository.DeleteByAccount(1)

	assert.NoError(t, err)
}
