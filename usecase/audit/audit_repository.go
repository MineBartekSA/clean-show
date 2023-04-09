package audit

import (
	"github.com/minebarteksa/clean-show/domain"
)

type auditRepository struct {
	db domain.DB

	insert domain.Stmt
}

func NewAuditRepository(db domain.DB) domain.AuditRepository {
	return &auditRepository{
		db:     db,
		insert: db.PrepareInsertStruct("audit_log", &domain.AuditEntry{}),
	}
}

func (ar *auditRepository) Insert(entry domain.AuditEntry) error {
	_, err := ar.insert.Exec(&entry)
	return domain.SQLError(err)
}

func (ar *auditRepository) BatchInsert(entries []domain.AuditEntry) error {
	for _, entry := range entries {
		_, err := ar.insert.Exec(&entry)
		if err != nil {
			return domain.SQLError(err)
		}
	}
	return nil
}
