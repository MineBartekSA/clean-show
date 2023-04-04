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
	_, err := ar.insert.Exec(ar.db.PrepareStruct(&entry))
	return err
}
