package order_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/logger"
	"github.com/minebarteksa/clean-show/test"
	"github.com/minebarteksa/clean-show/usecase/order"
	"github.com/stretchr/testify/assert"
)

var (
	repositoryCache domain.OrderRepository
	mockCache       sqlmock.Sqlmock
	preparedCache   []*sqlmock.ExpectedPrepare
)

func NewRepository(t *testing.T) (domain.OrderRepository, sqlmock.Sqlmock, []*sqlmock.ExpectedPrepare) {
	if repositoryCache == nil {
		logger.InitDebug()
		test.SetupConfig()
		db, mock := test.NewMockDB(t)
		mockCache = mock

		preparedCache = []*sqlmock.ExpectedPrepare{
			mock.ExpectPrepare("SELECT COUNT\\(\\*\\) FROM orders"),
			mock.ExpectPrepare("SELECT \\* FROM orders WHERE deleted_at IS NULL ORDER BY id DESC LIMIT \\? OFFSET \\?"),
			mock.ExpectPrepare("SELECT .* FROM orders WHERE order_by = \\? AND deleted_at IS NULL"),
			mock.ExpectPrepare("SELECT .* FROM orders WHERE id = \\? AND deleted_at IS NULL"),
			mock.ExpectPrepare("SELECT order_by FROM orders WHERE id = \\? AND deleted_at IS NULL"),
			mock.ExpectPrepare("INSERT INTO orders \\(.*\\) VALUES \\(.*\\) RETURNING id"),
			mock.ExpectPrepare("UPDATE orders SET .*, updated_at = NOW\\(\\) WHERE id = \\?"),
			mock.ExpectPrepare("UPDATE orders SET status = \\?, updated_at = NOW\\(\\) WHERE id = \\?"),
			mock.ExpectPrepare("UPDATE orders SET deleted_at = NOW\\(\\) WHERE id = \\?"),
		}

		repositoryCache = order.NewOrderRepository(db)
	}
	return repositoryCache, mockCache, preparedCache
}

func TestCount(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	prepared[0].ExpectQuery().WillReturnRows(test.NewRows("COUNT(*)").AddRow(3))

	count, err := repository.Count()

	assert.NoError(t, err)
	assert.Equal(t, uint(3), count)
}

var mockOrders = []domain.Order{
	{
		DBModel: domain.DBModel{
			ID: 1,
		},
		Status:  domain.OrderStatusInRealisation,
		OrderBy: 1,
		Products: domain.DBArray[domain.ProductOrder]{
			{
				ProductID: 1,
				Amount:    1,
				Price:     5,
			},
			{
				ProductID: 2,
				Amount:    1,
				Price:     6,
			},
		},
	},
	{
		DBModel: domain.DBModel{
			ID: 2,
		},
		Status:  domain.OrderStatusCreated,
		OrderBy: 1,
		Products: domain.DBArray[domain.ProductOrder]{
			{
				ProductID: 1,
				Amount:    1,
				Price:     5,
			},
		},
	},
	{
		DBModel: domain.DBModel{
			ID: 3,
		},
		Status:  domain.OrderStatusCompleted,
		OrderBy: 1,
		Products: domain.DBArray[domain.ProductOrder]{
			{
				ProductID: 1,
				Amount:    1,
				Price:     5,
			},
			{
				ProductID: 2,
				Amount:    1,
				Price:     6,
			},
			{
				ProductID: 3,
				Amount:    10,
				Price:     5,
			},
		},
	},
}

func TestSelect(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	rows := test.NewRows("id", "status", "order_by", "shipping_address", "invoice_address", "products", "shipping_price", "total")
	for _, order := range mockOrders {
		rows.AddRow(order.ID, order.Status, order.OrderBy, order.ShippingAddress, order.InvoiceAddress, order.Products, order.ShippingPrice, order.Total)
	}
	prepared[1].ExpectQuery().WillReturnRows(rows)

	orders, err := repository.Select(5, 2)

	assert.NoError(t, err)
	assert.NotEmpty(t, orders)
	assert.Len(t, orders[2].Products, 3)
	assert.Equal(t, mockOrders, orders)
}

func TestSelectAccount(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	rows := test.NewRows("id", "status", "order_by", "shipping_address", "invoice_address", "products", "shipping_price", "total")
	for _, order := range mockOrders {
		rows.AddRow(order.ID, order.Status, order.OrderBy, order.ShippingAddress, order.InvoiceAddress, order.Products, order.ShippingPrice, order.Total)
	}
	prepared[2].ExpectQuery().WithArgs(7).WillReturnRows(rows)

	orders, err := repository.SelectAccount(7)

	assert.NoError(t, err)
	assert.Equal(t, mockOrders, orders)
}

func TestSelectID(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	rows := test.NewRows("id", "status", "order_by", "shipping_address", "invoice_address", "products", "shipping_price", "total")
	mockOrder := mockOrders[2]
	rows.AddRow(mockOrder.ID, mockOrder.Status, mockOrder.OrderBy, mockOrder.ShippingAddress, mockOrder.InvoiceAddress, mockOrder.Products, mockOrder.ShippingPrice, mockOrder.Total)
	prepared[3].ExpectQuery().WithArgs(3).WillReturnRows(rows)

	order, err := repository.SelectID(3)

	assert.NoError(t, err)
	assert.Equal(t, mockOrders[2], *order)
}

func TestSelectOrderBy(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	prepared[4].ExpectQuery().WithArgs(3).WillReturnRows(test.NewRows("order_by").AddRow(1))

	ordered_by, err := repository.SelectOrderBy(3)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), ordered_by)
}

func TestInsert(t *testing.T) {
	repository, _, prepared := NewRepository(t)
	order := mockOrders[1]

	prepared[5].ExpectQuery().
		WithArgs(order.Status, order.OrderBy, order.ShippingAddress, order.InvoiceAddress, order.Products, order.ShippingPrice, order.Total).
		WillReturnRows(test.IDRow(1))

	err := repository.Insert(&order)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), order.ID)
}

func TestUpdate(t *testing.T) {
	repository, _, prepared := NewRepository(t)
	order := mockOrders[0]

	prepared[6].ExpectExec().
		WithArgs(order.Status, order.OrderBy, order.ShippingAddress, order.InvoiceAddress, order.Products, order.ShippingPrice, order.Total, order.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repository.Update(&order)

	assert.NoError(t, err)
}

func TestUpdateStatus(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	prepared[7].ExpectExec().WithArgs(domain.OrderStatusCompleted, 3).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repository.UpdateStatus(3, domain.OrderStatusCompleted)

	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	repository, _, prepared := NewRepository(t)

	prepared[8].ExpectExec().WithArgs(3).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repository.Delete(3)

	assert.NoError(t, err)
}
