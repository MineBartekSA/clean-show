package audit

import (
	"log"

	"github.com/minebarteksa/clean-show/domain"
)

type auditRepository struct {
	db domain.DB

	insert domain.Stmt
}

func NewAuditRepository(db domain.DB) domain.AuditRepository {
	insert, err := db.PrepareInsertStruct("audit_log", &domain.AuditEntry{})
	if err != nil {
		log.Panicf("failed to prepare a named insert statement from a structure: %s", err)
	}
	return &auditRepository{db, insert}
}

func (ar *auditRepository) Insert(entry domain.AuditEntry) error {
	_, err := ar.insert.Exec(ar.db.PrepareStruct(&entry))
	return err
}
