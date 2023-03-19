package product

import "github.com/minebarteksa/clean-show/domain"

type productUsecase struct {
	repository domain.ProductRepository
}

func NewProductUsecase(repository domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{repository}
}

func (u *productUsecase) ID(id uint) (*domain.Product, error) {
	return u.repository.FetchByID(id)
}
