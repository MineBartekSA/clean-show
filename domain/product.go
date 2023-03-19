package domain

type Product struct {
	//
}

type ProductController interface {
	Register(router Router)
	GetByID(context Context, session Session)
}

type ProductUsecase interface {
	ID(id uint) (*Product, error)
}

type ProductRepository interface {
	FetchByID(id uint) (*Product, error)
}
