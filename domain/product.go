package domain

type Product struct {
	DBModel     `json:"-"`
	Status      ProductStatus   `db:"status"`
	Name        string          `db:"name"`
	Description string          `db:"description"`
	Price       float64         `db:"price"`
	Images      DBArray[string] `db:"images"`
}

type ProductStatus int

const (
	ProductStatusInStock ProductStatus = iota + 1
	ProductStatusOutOfStock
	ProductStatusDiscontinued
)

type ProductController interface {
	Register(router Router)
	CreateNew(context Context, session UserSession)
	GetByID(context Context, session UserSession)
}

type ProductUsecase interface {
	Create(accountId uint, product *Product) error
	FetchByID(id uint) (*Product, error)
}

type ProductRepository interface {
	ID(id uint) (*Product, error)
	Insert(product *Product) error
}
