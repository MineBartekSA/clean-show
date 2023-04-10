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
	usecase := mocks.NewProductUsecase(t)
	controller := product.NewProductController(usecase)

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

	usecase.On("TotalCount").Return(uint(7), nil)
	usecase.On("Fetch", 5, 2).Return(products, nil)

	controller.Get(context, nil)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	data := context.Out.(domain.DataList[domain.Product])
	assert.Equal(t, uint(7), data.Hits)
	assert.Equal(t, uint(2), data.Pages)
	assert.Equal(t, products, data.Data)
}

func TestPost(t *testing.T) {
	usecase := mocks.NewProductUsecase(t)
	controller := product.NewProductController(usecase)

	session := test.NewUserSession(8)
	product := domain.Product{}
	context := &test.MockContext{
		In: &product,
	}
	usecase.On("Create", session.GetAccountID(), &product).Return(nil)

	controller.Post(context, session)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	data := context.Out.(struct {
		ID uint `json:"id"`
		*domain.Product
	})
	assert.Equal(t, uint(0), data.ID)
	assert.Equal(t, &product, data.Product)
}

func TestGetByID(t *testing.T) {
	usecase := mocks.NewProductUsecase(t)
	controller := product.NewProductController(usecase)

	product := domain.Product{
		DBModel: domain.DBModel{ID: 9},
	}
	context := &test.MockContext{
		ParamMap: map[string]string{
			"id": "9",
		},
	}

	usecase.On("FetchByID", uint(9)).Return(&product, nil)

	controller.GetByID(context, nil)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	data := context.Out.(*domain.Product)
	assert.Equal(t, &product, data)
}

func TestPatch(t *testing.T) {
	usecase := mocks.NewProductUsecase(t)
	controller := product.NewProductController(usecase)

	product := domain.Product{DBModel: domain.DBModel{ID: 8}}
	session := test.NewUserSession(1)
	context := &test.MockContext{
		ParamMap: map[string]string{
			"id": "8",
		},
		In: map[string]any{},
	}

	usecase.On("Modify", uint(1), uint(8), context.In).Return(&product, nil)

	controller.Patch(context, session)

	assert.Equal(t, http.StatusOK, context.OutStatus)
	data := context.Out.(struct {
		ID uint `json:"id"`
		*domain.Product
	})
	assert.Equal(t, uint(8), data.ID)
	assert.Equal(t, &product, data.Product)
}

func TestDeleteController(t *testing.T) {
	usecase := mocks.NewProductUsecase(t)
	controller := product.NewProductController(usecase)

	session := test.NewUserSession(7)
	context := &test.MockContext{
		ParamMap: map[string]string{
			"id": "1",
		},
	}

	usecase.On("Remove", uint(7), uint(1)).Return(nil)

	controller.Delete(context, session)

	assert.Equal(t, http.StatusNoContent, context.OutStatus)
}
