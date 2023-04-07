package product_test

import (
	"net/http"
	"testing"

	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/domain/mocks"
	"github.com/minebarteksa/clean-show/test"
	"github.com/minebarteksa/clean-show/usecase/product"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	mockUsecase := mocks.NewProductUsecase(t)
	controller := product.NewProductController(mockUsecase)

	context := &test.MockContext{
		QueryMap: map[string]string{
			"limit": "5",
			"page":  "2",
		},
	}
	products := []domain.Product{
		{Name: "TEST"},
		{Name: "TEST2"},
	}

	mockUsecase.On("TotalCount").Return(uint(7), nil)
	mockUsecase.On("Fetch", 5, 2).Return(products, nil)

	controller.Get(context, nil)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	data, ok := context.Out.(domain.DataList[domain.Product])
	assert.True(t, ok)
	if ok {
		assert.Equal(t, uint(7), data.Hits)
		assert.Equal(t, uint(2), data.Pages)
		assert.Equal(t, products, data.Data)
	}
}

// TODO: Write
