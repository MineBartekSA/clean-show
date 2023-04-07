package product_test

import (
	"testing"

	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/domain/mocks"
	"github.com/minebarteksa/clean-show/test"
	"github.com/minebarteksa/clean-show/usecase/product"
	"github.com/stretchr/testify/assert"
)

func TestTotalCount(t *testing.T) {
	repository := mocks.NewProductRepository(t)
	_, audit := test.NewAuditUsecase(t, domain.ResourceTypeProduct)
	usecase := product.NewProductUsecase(repository, audit)

	repository.On("Count").Return(uint(7), nil)

	count, err := usecase.TotalCount()

	assert.NoError(t, err)
	assert.Equal(t, uint(7), count)
}

func TestFetch(t *testing.T) {
	repository := mocks.NewProductRepository(t)
	_, audit := test.NewAuditUsecase(t, domain.ResourceTypeProduct)
	usecase := product.NewProductUsecase(repository, audit)

	repository.On("Select", 5, 2).Return([]domain.Product{{}, {}}, nil)

	products, err := usecase.Fetch(5, 2)

	assert.NoError(t, err)
	assert.Len(t, products, 2)
}

func TestCreate(t *testing.T) {
	repository := mocks.NewProductRepository(t)
	resource, audit := test.NewAuditUsecase(t, domain.ResourceTypeProduct)
	usecase := product.NewProductUsecase(repository, audit)

	product := domain.Product{DBModel: domain.DBModel{ID: 1}}
	repository.On("Insert", &product).Return(nil)
	resource.On("Creation", uint(1), uint(1)).Return(nil)

	err := usecase.Create(1, &product)

	assert.NoError(t, err)
}

func TestModify(t *testing.T) {
	repository := mocks.NewProductRepository(t)
	resource, audit := test.NewAuditUsecase(t, domain.ResourceTypeProduct)
	usecase := product.NewProductUsecase(repository, audit)

	product := domain.Product{}
	repository.On("SelectID", uint(1)).Return(&product, nil)
	repository.On("Update", &product).Return(nil)
	resource.On("Modification", uint(1), uint(1)).Return(nil)

	err := usecase.Modify(1, 1, map[string]any{"name": "hello"})

	assert.NoError(t, err)
	assert.Equal(t, "hello", product.Name)
}

func TestRemove(t *testing.T) {
	repository := mocks.NewProductRepository(t)
	resource, audit := test.NewAuditUsecase(t, domain.ResourceTypeProduct)
	usecase := product.NewProductUsecase(repository, audit)

	repository.On("Delete", uint(1)).Return(nil)
	resource.On("Deletion", uint(1), uint(1)).Return(nil)

	err := usecase.Remove(1, 1)

	assert.NoError(t, err)
}
