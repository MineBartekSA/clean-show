package product

import "github.com/minebarteksa/clean-show/domain"

type productUsecase struct {
	repository domain.ProductRepository
	audit      domain.AuditResource
}

func NewProductUsecase(repository domain.ProductRepository, audit domain.AuditUsecase) domain.ProductUsecase {
	return &productUsecase{repository, audit.Resource(domain.ResourceTypeProduct)}
}

func (pu *productUsecase) Create(accountId uint, product *domain.Product) error {
	err := pu.repository.Insert(product)
	if err != nil {
		return err
	}
	return pu.audit.Creation(accountId, product.ID)
}

func (pu *productUsecase) FetchByID(id uint) (*domain.Product, error) {
	return pu.repository.ID(id)
}
