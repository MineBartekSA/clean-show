package product

import (
	"fmt"

	"github.com/minebarteksa/clean-show/domain"
)

type productRepository struct {
	db domain.DB
}

func NewProductRepository(db domain.DB) domain.ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) FetchByID(id uint) (*domain.Product, error) {
	// TODO: Write
	return nil, fmt.Errorf("not implemented")
}
