package domain

type Product struct {
	DBModel     `json:"-"`
	Status      ProductStatus   `db:"status" json:"status"`
	Name        string          `db:"name" json:"name"`
	Description string          `db:"description" json:"description"`
	Price       float64         `db:"price" json:"price"`
	Images      DBArray[string] `db:"images" json:"images"`
}

type ProductStatus int

const (
	ProductStatusInStock ProductStatus = iota + 1
	ProductStatusOutOfStock
	ProductStatusDiscontinued
)

//go:generate mockery --name ProductController
type ProductController interface {
	Register(router Router)
	Get(context Context, session UserSession)
	Post(context Context, session UserSession)
	GetByID(context Context, session UserSession)
	Patch(context Context, session UserSession)
	Delete(context Context, session UserSession)
}

//go:generate mockery --name ProductUsecase
type ProductUsecase interface {
	TotalCount() (uint, error)
	Fetch(limit, page int) ([]Product, error)
	Create(accountId uint, product *Product) error
	FetchByID(id uint) (*Product, error)
	Modify(accountId, productId uint, data map[string]any) (*Product, error)
	Remove(accountId, productId uint) error
}

//go:generate mockery --name ProductRepository
type ProductRepository interface {
	Count() (uint, error)
	Select(limit, page int) ([]Product, error)
	SelectID(id uint) (*Product, error)
	Insert(product *Product) error
	Update(product *Product) error
	Delete(id uint) error
}
