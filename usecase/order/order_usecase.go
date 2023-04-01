package order

import (
	"fmt"

	"github.com/minebarteksa/clean-show/domain"
)

type orderUsecase struct {
	repository domain.OrderRepository
	audit      domain.AuditResource
}

func NewOrderUsecase(repository domain.OrderRepository, audit domain.AuditUsecase) domain.OrderUsecase {
	return &orderUsecase{repository, audit.Resource(domain.ResourceTypeOrder)}
}

func (ou *orderUsecase) TotalCount() (uint, error) {
	//
	return 0, fmt.Errorf("not implemented")
}

func (ou *orderUsecase) Fetch(limit, page int) ([]domain.Order, error) {
	//
	return nil, fmt.Errorf("not implemented")
}

func (ou *orderUsecase) Create(order *domain.Order) error {
	//
	return fmt.Errorf("not implemented")
}

func (ou *orderUsecase) FetchByID(session domain.UserSession, id uint) (*domain.Order, error) {
	//
	return nil, fmt.Errorf("not implemented")
}

func (ou *orderUsecase) Modify(accountId, orderId uint, data map[string]any) error {
	//
	return fmt.Errorf("not implemented")
}

func (ou *orderUsecase) Cancel(session domain.UserSession, orderId uint) error {
	//
	return fmt.Errorf("not implemented")
}

func (ou *orderUsecase) Delete(accountId, orderId uint) error {
	//
	return fmt.Errorf("not implemented")
}
