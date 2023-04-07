package product_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/logger"
	"github.com/minebarteksa/clean-show/test"
	"github.com/minebarteksa/clean-show/usecase/product"
	"github.com/stretchr/testify/assert"
)

var (
	repositoryCache domain.ProductRepository
	mockCache       sqlmock.Sqlmock
	preparedCache   []*sqlmock.ExpectedPrepare
)

func NewRepository(t *testing.T) (domain.ProductRepository, sqlmock.Sqlmock, []*sqlmock.ExpectedPrepare) {
	if repositoryCache == nil {
		logger.InitDebug()
		test.SetupConfig()
		db, mock := test.NewMockDB(t)
		mockCache = mock

		preparedCache = []*sqlmock.ExpectedPrepare{
			mock.ExpectPrepare("SELECT COUNT\\(\\*\\) FROM products"),
			mock.ExpectPrepare("SELECT \\* FROM products WHERE deleted_at IS NULL LIMIT \\? OFFSET \\?"),
			mock.ExpectPrepare("SELECT .* FROM products WHERE id = \\? AND deleted_at IS NULL"),
			mock.ExpectPrepare("INSERT INTO products \\(.*\\) VALUES \\(.*\\) RETURNING id"),
			mock.ExpectPrepare("UPDATE products SET .*, updated_at = NOW\\(\\) WHERE id = \\?"),
			mock.ExpectPrepare("UPDATE products SET deleted_at = NOW\\(\\) WHERE id = \\?"),
		}

		repositoryCache = product.NewProductRepository(db)
	}
	return repositoryCache, mockCache, preparedCache
}

func TestCount(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	prepared[0].ExpectQuery().WillReturnRows(test.NewRows("COUNT(*)").AddRow(7))

	count, err := repository.Count()

	assert.NoError(t, err)
	assert.Equal(t, uint(7), count)
}

func TestSelect(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	prepared[1].ExpectQuery().WithArgs(5, 5).WillReturnRows(
		test.NewRows("id", "status", "name", "description", "price", "images").
			AddRow(1, 1, "test 1", "", float64(1), "test;hello;world").
			AddRow(2, 1, "test 2", "", float64(2), "world;hello;test"),
	)

	products, err := repository.Select(5, 2)

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, domain.Product{
		DBModel: domain.DBModel{
			ID: 1,
		},
		Status: domain.ProductStatusInStock,
		Name:   "test 1",
		Price:  1,
		Images: domain.DBArray[string]{
			"test",
			"hello",
			"world",
		},
	}, products[0])
	assert.Equal(t, domain.Product{
		DBModel: domain.DBModel{
			ID: 2,
		},
		Status: domain.ProductStatusInStock,
		Name:   "test 2",
		Price:  2,
		Images: domain.DBArray[string]{
			"world",
			"hello",
			"test",
		},
	}, products[1])
}

func TestSelectID(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	prepared[2].ExpectQuery().WithArgs(1).WillReturnRows(
		test.NewRows("id", "status", "name", "description", "price", "images").
			AddRow(1, 1, "test 1", "", float64(1), "test;hello;world"),
	)

	product, err := repository.SelectID(1)

	assert.NoError(t, err)
	assert.Equal(t, domain.Product{
		DBModel: domain.DBModel{
			ID: 1,
		},
		Status: domain.ProductStatusInStock,
		Name:   "test 1",
		Price:  1,
		Images: domain.DBArray[string]{
			"test",
			"hello",
			"world",
		},
	}, *product)
}

func TestInsert(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	product := domain.Product{
		Status: domain.ProductStatusOutOfStock,
		Name:   "test 1",
		Price:  2,
		Images: domain.DBArray[string]{
			"tt",
			"hh",
		},
	}

	prepared[3].ExpectQuery().
		WithArgs(product.Status, product.Name, product.Description, product.Price, product.Images).
		WillReturnRows(test.NewRows("id").AddRow(1))

	err := repository.Insert(&product)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), product.ID)
}

func TestUpdate(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	product := domain.Product{
		DBModel: domain.DBModel{
			ID: 1,
		},
		Status: domain.ProductStatusOutOfStock,
		Name:   "test 1",
		Price:  2,
		Images: domain.DBArray[string]{
			"tt",
			"hh",
		},
	}

	prepared[4].ExpectExec().
		WithArgs(product.Status, product.Name, product.Description, product.Price, product.Images, product.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repository.Update(&product)

	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	prepared[5].ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repository.Delete(1)

	assert.NoError(t, err)
}
