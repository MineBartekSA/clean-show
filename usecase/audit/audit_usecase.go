package audit

import (
	"github.com/minebarteksa/clean-show/domain"
)

type auditUsecase struct {
	repository domain.AuditRepository
}

func NewAuditUsecase(repository domain.AuditRepository) domain.AuditUsecase {
	return &auditUsecase{repository}
}

func (au *auditUsecase) Create(entry_type domain.EntryType, resource_type domain.ResourceType, resource_id uint, executor uint) error {
	return au.repository.Insert(domain.AuditEntry{
		Type:         entry_type,
		ResourceType: resource_type,
		ResourceID:   resource_id,
		ExecutorID:   executor,
	}) // TODO: Encapsulate error?
}

func (au *auditUsecase) Creation(executor uint, res_type domain.ResourceType, res_id uint) error {
	return au.Create(domain.EntryTypeCreation, res_type, res_id, executor)
}

func (au *auditUsecase) Modification(executor uint, res_type domain.ResourceType, res_id uint) error {
	return au.Create(domain.EntryTypeModification, res_type, res_id, executor)
}

func (au *auditUsecase) Deletion(executor uint, res_type domain.ResourceType, res_id uint) error {
	return au.Create(domain.EntryTypeDeletion, res_type, res_id, executor)
}

func (au *auditUsecase) Resource(resource_type domain.ResourceType) domain.AuditResource {
	return &auditResource{au, resource_type}
}

type auditResource struct {
	usecase  domain.AuditUsecase
	resource domain.ResourceType
}

func (ar *auditResource) Creation(executor uint, resource_id uint) error {
	return ar.usecase.Creation(executor, ar.resource, resource_id)
}

func (ar *auditResource) Modification(executor uint, resource_id uint) error {
	return ar.usecase.Modification(executor, ar.resource, resource_id)
}

func (ar *auditResource) Deletion(executor uint, resource_id uint) error {
	return ar.usecase.Deletion(executor, ar.resource, resource_id)
}
