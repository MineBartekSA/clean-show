package domain

type AuditEntry struct {
	DBModel
	Type         EntryType    `db:"type"`
	ResourceType ResourceType `db:"resource_type"`
	ResourceID   uint         `db:"resource_id"`
	ExecutorID   uint         `db:"executor_id"`
}

type EntryType int

const (
	EntryTypeCreation EntryType = iota + 1
	EntryTypeModification
	EntryTypeDeletion
)

type ResourceType int

const (
	ResourceTypeProduct ResourceType = iota + 1
	ResourceTypeOrder
	ResourceTypeAccount
	ResourceTypeAccountPassword
	ResourceTypeSession
)

//go:generate mockery --name AuditResource
type AuditResource interface {
	Creation(executor uint, resId uint) error
	Modification(executor uint, resId uint) error
	BatchModification(executor uint, resIds []uint) error
	Deletion(executor uint, resId uint) error
}

//go:generate mockery --name AuditUsecase
type AuditUsecase interface {
	Create(entry_type EntryType, resource_type ResourceType, resourceId uint, executor uint) error

	Creation(executor uint, resType ResourceType, resId uint) error
	Modification(executor uint, resType ResourceType, resid uint) error
	BatchModification(executor uint, resType ResourceType, resIds []uint) error
	Deletion(executor uint, resType ResourceType, resId uint) error

	Resource(resource_type ResourceType) AuditResource
}

//go:generate mockery --name AuditRepository
type AuditRepository interface {
	Insert(entry AuditEntry) error
	BatchInsert(entry []AuditEntry) error
}
