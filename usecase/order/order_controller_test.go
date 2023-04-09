package order_test

import (
	"net/http"
	"testing"

	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/domain/mocks"
	"github.com/minebarteksa/clean-show/test"
	"github.com/minebarteksa/clean-show/usecase/order"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	usecase := mocks.NewOrderUsecase(t)
	controller := order.NewOrderController(usecase)

	context := &test.MockContext{
		QueryMap: map[string]string{
			"limit": "5",
			"page":  "2",
		},
	}
	orders := []domain.Order{{Status: domain.OrderStatusCanceled}}

	usecase.On("TotalCount").Return(uint(10), nil)
	usecase.On("Fetch", 5, 2).Return(orders, nil)

	controller.Get(context, nil)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	data := context.Out.(domain.DataList[domain.Order])
	assert.Equal(t, uint(10), data.Hits)
	assert.Equal(t, uint(2), data.Pages)
	assert.Equal(t, orders, data.Data)
}

func TestPost(t *testing.T) {
	usecase := mocks.NewOrderUsecase(t)
	controller := order.NewOrderController(usecase)

	create := domain.OrderCreate{}
	session := test.NewUserSession(1)
	context := &test.MockContext{
		In: &create,
	}

	order := create.ToOrder(1)
	order.ID = 10
	usecase.On("Create", uint(1), context.In).Return(order, nil)

	controller.Post(context, session)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	data := context.Out.(struct {
		ID uint `json:"id"`
		*domain.Order
	})
	assert.Equal(t, uint(10), data.ID)
	assert.Equal(t, order, data.Order)
}

func TestGetByID(t *testing.T) {
	usecase := mocks.NewOrderUsecase(t)
	controller := order.NewOrderController(usecase)

	order := domain.Order{}
	session := test.NewUserSession(1)
	context := &test.MockContext{
		ParamMap: map[string]string{
			"id": "10",
		},
	}

	usecase.On("FetchByID", session, uint(10)).Return(&order, nil)

	controller.GetByID(context, session)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	data := context.Out.(*domain.Order)
	assert.Equal(t, &order, data)
}

func TestPatch(t *testing.T) {
	usecase := mocks.NewOrderUsecase(t)
	controller := order.NewOrderController(usecase)

	session := test.NewUserSession(1)
	context := &test.MockContext{
		ParamMap: map[string]string{
			"id": "10",
		},
		In: map[string]any{},
	}

	usecase.On("Modify", uint(1), uint(10), context.In).Return(nil)

	controller.Patch(context, session)

	assert.Equal(t, http.StatusNoContent, context.OutStatus)
}

func TestPostCancel(t *testing.T) {
	usecase := mocks.NewOrderUsecase(t)
	controller := order.NewOrderController(usecase)

	session := test.NewUserSession(1)
	context := &test.MockContext{
		ParamMap: map[string]string{
			"id": "10",
		},
	}

	usecase.On("Cancel", session, uint(10)).Return(nil)

	controller.PostCancel(context, session)

	assert.Equal(t, http.StatusNoContent, context.OutStatus)
}

func TestDeleteController(t *testing.T) {
	usecase := mocks.NewOrderUsecase(t)
	controller := order.NewOrderController(usecase)

	session := test.NewUserSession(1)
	context := &test.MockContext{
		ParamMap: map[string]string{
			"id": "10",
		},
	}

	usecase.On("Remove", uint(1), uint(10)).Return(nil)

	controller.Delete(context, session)

	assert.Equal(t, http.StatusNoContent, context.OutStatus)
}
