package product

import (
	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/usecase"
)

type productUsecase struct {
	repository domain.ProductRepository
	audit      domain.AuditResource
}

func NewProductUsecase(repository domain.ProductRepository, audit domain.AuditUsecase) domain.ProductUsecase {
	return &productUsecase{repository, audit.Resource(domain.ResourceTypeProduct)}
}

func (pu *productUsecase) TotalCount() (uint, error) {
	return pu.repository.Count()
}

func (pu *productUsecase) Fetch(limit, page int) ([]domain.Product, error) {
	if limit < 0 {
		limit = 0
	} else if limit > 1000 {
		limit = 1000
	}
	if page < 1 {
		page = 1
	}
	list, err := pu.repository.Select(limit, page)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (pu *productUsecase) Create(accountId uint, product *domain.Product) error {
	err := pu.repository.Insert(product)
	if err != nil {
		return err
	}
	return pu.audit.Creation(accountId, product.ID)
}

func (pu *productUsecase) FetchByID(id uint) (*domain.Product, error) {
	return pu.repository.SelectID(id)
}

func (pu *productUsecase) Modify(accountId, productId uint, data map[string]any) (*domain.Product, error) {
	product, err := pu.repository.SelectID(productId)
	if err != nil {
		return nil, err
	}
	err = usecase.PatchModel(product, data)
	if err != nil {
		return nil, err
	}
	err = pu.repository.Update(product)
	if err != nil {
		return nil, err
	}
	return product, pu.audit.Modification(accountId, productId)
}

func (pu *productUsecase) Remove(accountId, productId uint) error {
	err := pu.repository.Delete(productId)
	if err != nil {
		return err
	}
	return pu.audit.Deletion(accountId, productId)
}
