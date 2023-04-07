package order_test

import (
	"testing"

	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/domain/mocks"
	"github.com/minebarteksa/clean-show/test"
	"github.com/minebarteksa/clean-show/usecase/order"
	"github.com/stretchr/testify/assert"
)

func TestTotalCount(t *testing.T) {
	repository := mocks.NewOrderRepository(t)
	_, audit := test.NewAuditUsecase(t, domain.ResourceTypeOrder)
	usecase := order.NewOrderUsecase(repository, audit)

	repository.On("Count").Return(uint(10), nil)

	count, err := usecase.TotalCount()

	assert.NoError(t, err)
	assert.Equal(t, uint(10), count)
}

func TestFetch(t *testing.T) {
	repository := mocks.NewOrderRepository(t)
	_, audit := test.NewAuditUsecase(t, domain.ResourceTypeOrder)
	usecase := order.NewOrderUsecase(repository, audit)

	repository.On("Select", 5, 2).Return([]domain.Order{{}, {}}, nil)

	orders, err := usecase.Fetch(5, 2)

	assert.NoError(t, err)
	assert.Len(t, orders, 2)
}

func TestFetchByAccount(t *testing.T) {
	repository := mocks.NewOrderRepository(t)
	_, audit := test.NewAuditUsecase(t, domain.ResourceTypeOrder)
	usecase := order.NewOrderUsecase(repository, audit)

	repository.On("SelectAccount", uint(1)).Return([]domain.Order{{}, {}, {}}, nil)

	orders, err := usecase.FetchByAccount(1)

	assert.NoError(t, err)
	assert.Len(t, orders, 3)
}

func TestCreate(t *testing.T) {
	repository := mocks.NewOrderRepository(t)
	resource, audit := test.NewAuditUsecase(t, domain.ResourceTypeOrder)
	usecase := order.NewOrderUsecase(repository, audit)

	create := domain.OrderCreate{}
	mockOrder := create.ToOrder(1)

	repository.On("Insert", mockOrder).Return(nil)
	resource.On("Creation", uint(1), uint(0)).Return(nil)

	order, err := usecase.Create(1, &create)

	assert.NoError(t, err)
	assert.Equal(t, mockOrder, order)
}

func TestFetchByID(t *testing.T) {
	repository := mocks.NewOrderRepository(t)
	_, audit := test.NewAuditUsecase(t, domain.ResourceTypeOrder)
	usecase := order.NewOrderUsecase(repository, audit)

	session := test.MockUserSession{
		Account: &domain.Account{
			DBModel: domain.DBModel{
				ID: 1,
			},
		},
	}
	order := domain.Order{
		DBModel: domain.DBModel{
			ID: 10,
		},
		OrderBy: 1,
	}

	repository.On("SelectID", uint(1)).Return(&order, nil)

	out, err := usecase.FetchByID(&session, 1)

	assert.NoError(t, err)
	assert.Equal(t, order, *out)
}

func TestModify(t *testing.T) {
	repository := mocks.NewOrderRepository(t)
	resource, audit := test.NewAuditUsecase(t, domain.ResourceTypeOrder)
	usecase := order.NewOrderUsecase(repository, audit)

	order := domain.Order{
		DBModel: domain.DBModel{
			ID: 7,
		},
	}
	session := test.MockUserSession{
		Account: &domain.Account{
			DBModel: domain.DBModel{
				ID: 1,
			},
		},
	}

	repository.On("SelectID", order.ID).Return(&order, nil)
	repository.On("Update", &order).Return(nil)
	resource.On("Modification", session.GetAccountID(), order.ID).Return(nil)

	err := usecase.Modify(uint(1), 7, map[string]any{
		"shipping_price": float64(10.5),
		"total":          float64(1000),
		"status":         domain.OrderStatusShipped,
	})

	assert.NoError(t, err)
	assert.Equal(t, domain.OrderStatusShipped, order.Status)
	assert.Equal(t, float64(10.5), order.ShippingPrice)
	assert.Empty(t, order.Total)
}

func TestCancel(t *testing.T) {
	repository := mocks.NewOrderRepository(t)
	resource, audit := test.NewAuditUsecase(t, domain.ResourceTypeOrder)
	usecase := order.NewOrderUsecase(repository, audit)

	session := test.MockUserSession{
		Account: &domain.Account{
			DBModel: domain.DBModel{
				ID: 1,
			},
		},
	}

	repository.On("SelectOrderBy", uint(3)).Return(uint(1), nil)
	repository.On("UpdateStatus", uint(3), domain.OrderStatusCanceled).Return(nil)
	resource.On("Modification", uint(1), uint(3)).Return(nil)

	err := usecase.Cancel(&session, 3)

	assert.NoError(t, err)
}

func TestCancelByAccount(t *testing.T) {
	repository := mocks.NewOrderRepository(t)
	resource, audit := test.NewAuditUsecase(t, domain.ResourceTypeOrder)
	usecase := order.NewOrderUsecase(repository, audit)

	repository.On("SelectAccount", uint(2)).Return([]domain.Order{
		{DBModel: domain.DBModel{ID: 1}},
		{DBModel: domain.DBModel{ID: 2}},
	}, nil)
	repository.On("UpdateStatus", uint(1), domain.OrderStatusCanceled).Return(nil)
	resource.On("Modification", uint(1), uint(1)).Return(nil)
	repository.On("UpdateStatus", uint(2), domain.OrderStatusCanceled).Return(nil)
	resource.On("Modification", uint(1), uint(2)).Return(nil)

	err := usecase.CancelByAccount(1, 2)

	assert.NoError(t, err)
}

func TestRemove(t *testing.T) {
	repository := mocks.NewOrderRepository(t)
	resource, audit := test.NewAuditUsecase(t, domain.ResourceTypeOrder)
	usecase := order.NewOrderUsecase(repository, audit)

	repository.On("Delete", uint(1)).Return(nil)
	resource.On("Deletion", uint(1), uint(1)).Return(nil)

	err := usecase.Remove(1, 1)

	assert.NoError(t, err)
}
