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
	data, ok := context.Out.(domain.DataList[domain.Order])
	assert.True(t, ok)
	if ok {
		assert.Equal(t, uint(10), data.Hits)
		assert.Equal(t, uint(2), data.Pages)
		assert.Equal(t, orders, data.Data)
	}
}

// TODO: Write
