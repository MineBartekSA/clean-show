package audit_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/logger"
	"github.com/minebarteksa/clean-show/test"
	"github.com/minebarteksa/clean-show/usecase/audit"
	"github.com/stretchr/testify/assert"
)

var (
	repositoryCache domain.AuditRepository
	mockCache       sqlmock.Sqlmock
	preparedCache   []*sqlmock.ExpectedPrepare
)

func NewRepository(t *testing.T) (domain.AuditRepository, sqlmock.Sqlmock, []*sqlmock.ExpectedPrepare) {
	if repositoryCache == nil {
		logger.InitDebug()
		test.SetupConfig()
		db, mock := test.NewMockDB(t)
		mockCache = mock

		preparedCache = []*sqlmock.ExpectedPrepare{
			mock.ExpectPrepare("INSERT INTO audit_log \\(.*\\) VALUES \\(.*\\) RETURNING id"),
		}

		repositoryCache = audit.NewAuditRepository(db)
	}
	return repositoryCache, mockCache, preparedCache
}

func TestInsert(t *testing.T) {
	repository, _, prepared := NewRepository(t)
	entry := domain.AuditEntry{
		Type:         domain.EntryTypeDeletion,
		ResourceType: domain.ResourceTypeSession,
		ResourceID:   10,
		ExecutorID:   50,
	}

	prepared[0].ExpectExec().
		WithArgs(entry.Type, entry.ResourceType, entry.ResourceID, entry.ExecutorID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.Insert(entry)

	assert.NoError(t, err)
}
