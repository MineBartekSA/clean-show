package audit_test

import (
	"testing"

	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/domain/mocks"
	"github.com/minebarteksa/clean-show/usecase/audit"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	repository := mocks.NewAuditRepository(t)
	usecase := audit.NewAuditUsecase(repository)

	entry := domain.AuditEntry{
		Type:         domain.EntryTypeModification,
		ResourceType: domain.ResourceTypeAccountPassword,
		ResourceID:   42,
		ExecutorID:   94,
	}

	repository.On("Insert", entry).Return(nil)

	err := usecase.Create(entry.Type, entry.ResourceType, entry.ResourceID, entry.ExecutorID)

	assert.NoError(t, err)
}
