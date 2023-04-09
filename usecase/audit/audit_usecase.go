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

func (au *auditUsecase) Create(entryType domain.EntryType, resourceType domain.ResourceType, resourceId uint, executor uint) error {
	return au.repository.Insert(domain.AuditEntry{
		Type:         entryType,
		ResourceType: resourceType,
		ResourceID:   resourceId,
		ExecutorID:   executor,
	})
}

func (au *auditUsecase) BatchCreate(entryType domain.EntryType, resourceType domain.ResourceType, resources []uint, executor uint) error {
	var entries []domain.AuditEntry
	for _, id := range resources {
		entries = append(entries, domain.AuditEntry{
			Type:         entryType,
			ResourceType: resourceType,
			ResourceID:   id,
			ExecutorID:   executor,
		})
	}
	return au.repository.BatchInsert(entries)
}

func (au *auditUsecase) Creation(executor uint, resType domain.ResourceType, resId uint) error {
	return au.Create(domain.EntryTypeCreation, resType, resId, executor)
}

func (au *auditUsecase) Modification(executor uint, resType domain.ResourceType, resId uint) error {
	return au.Create(domain.EntryTypeModification, resType, resId, executor)
}

func (au *auditUsecase) BatchModification(executor uint, resType domain.ResourceType, resIds []uint) error {
	return au.BatchCreate(domain.EntryTypeModification, resType, resIds, executor)
}

func (au *auditUsecase) Deletion(executor uint, resType domain.ResourceType, resId uint) error {
	return au.Create(domain.EntryTypeDeletion, resType, resId, executor)
}

func (au *auditUsecase) Resource(resource_type domain.ResourceType) domain.AuditResource {
	return &auditResource{au, resource_type}
}

type auditResource struct {
	usecase  domain.AuditUsecase
	resource domain.ResourceType
}

func (ar *auditResource) Creation(executor uint, resourceId uint) error {
	return ar.usecase.Creation(executor, ar.resource, resourceId)
}

func (ar *auditResource) Modification(executor uint, resourceId uint) error {
	return ar.usecase.Modification(executor, ar.resource, resourceId)
}

func (ar *auditResource) BatchModification(executor uint, resourceIds []uint) error {
	return ar.usecase.BatchModification(executor, ar.resource, resourceIds)
}

func (ar *auditResource) Deletion(executor uint, resourceId uint) error {
	return ar.usecase.Deletion(executor, ar.resource, resourceId)
}
