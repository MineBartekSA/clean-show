package order

import (
	"fmt"

	"github.com/minebarteksa/clean-show/domain"
)

type orderRepository struct {
	db domain.DB

	//
}

func NewOrderRepository(db domain.DB) domain.OrderRepository {
	return &orderRepository{db}
}

func (or *orderRepository) Count() (uint, error) {
	//
	return 0, fmt.Errorf("not implemented")
}

func (or *orderRepository) Select(limit, page int) ([]domain.Order, error) {
	//
	return nil, fmt.Errorf("not implemented")
}

func (or *orderRepository) SelectID(id uint) (*domain.Order, error) {
	//
	return nil, fmt.Errorf("not implemented")
}

func (or *orderRepository) Insert(order *domain.Order) error {
	//
	return fmt.Errorf("not implemented")
}

func (or *orderRepository) Update(order *domain.Order) error {
	//
	return fmt.Errorf("not implemented")
}

func (or *orderRepository) Delete(id uint) error {
	//
	return fmt.Errorf("not implemented")
}
