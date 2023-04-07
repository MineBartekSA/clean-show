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
	Creation(executor uint, res_id uint) error
	Modification(executor uint, res_id uint) error
	Deletion(executor uint, res_id uint) error
}

//go:generate mockery --name AuditUsecase
type AuditUsecase interface {
	Create(entry_type EntryType, resource_type ResourceType, resource_id uint, executor uint) error

	Creation(executor uint, res_type ResourceType, res_id uint) error
	Modification(executor uint, res_type ResourceType, res_id uint) error
	Deletion(executor uint, res_type ResourceType, res_id uint) error

	Resource(resource_type ResourceType) AuditResource
}

//go:generate mockery --name AuditRepository
type AuditRepository interface {
	Insert(entry AuditEntry) error
}
