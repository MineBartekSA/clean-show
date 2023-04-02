package order

import (
	"fmt"

	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/usecase"
)

type orderUsecase struct {
	repository domain.OrderRepository
	audit      domain.AuditResource
}

func NewOrderUsecase(repository domain.OrderRepository, audit domain.AuditUsecase) domain.OrderUsecase {
	return &orderUsecase{repository, audit.Resource(domain.ResourceTypeOrder)}
}

func (ou *orderUsecase) TotalCount() (uint, error) {
	return ou.repository.Count()
}

func (ou *orderUsecase) Fetch(limit, page int) ([]domain.Order, error) {
	if limit < 0 {
		limit = 0
	} else if limit > 1000 {
		limit = 1000
	}
	if page < 1 {
		page = 1
	}
	list, err := ou.repository.Select(limit, page)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (ou *orderUsecase) Create(accountId uint, create *domain.OrderCreate) (*domain.Order, error) {
	order := create.ToOrder(accountId)
	err := ou.repository.Insert(order)
	if err != nil {
		return nil, err
	}
	err = ou.audit.Creation(accountId, order.ID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (ou *orderUsecase) FetchByID(session domain.UserSession, id uint) (*domain.Order, error) {
	order, err := ou.repository.SelectID(id)
	if err != nil {
		return nil, err
	}
	if !session.IsStaff() && session.GetAccountID() != order.OrderBy {
		return nil, fmt.Errorf("only staff users can see other users orders")
	}
	return order, nil
}

func (ou *orderUsecase) Modify(accountId, orderId uint, data map[string]any) error {
	return ou.update(accountId, orderId, func(order *domain.Order) error {
		return usecase.PatchModel(order, data)
	})
}

func (ou *orderUsecase) Cancel(session domain.UserSession, orderId uint) error {
	aid := session.GetAccountID()
	return ou.update(aid, orderId, func(order *domain.Order) error {
		if !session.IsStaff() && aid != order.OrderBy {
			return fmt.Errorf("only staff users can cancel other users orders")
		}
		order.Status = domain.OrderStatusCanceled
		return nil
	})
}

func (ou *orderUsecase) Delete(accountId, orderId uint) error {
	err := ou.repository.Delete(orderId)
	if err != nil {
		return err
	}
	return ou.audit.Deletion(accountId, orderId)
}

func (ou *orderUsecase) update(accountId, orderId uint, mod func(order *domain.Order) error) error {
	order, err := ou.repository.SelectID(orderId)
	if err != nil {
		return err
	}
	err = mod(order)
	if err != nil {
		return err
	}
	err = ou.repository.Update(order)
	if err != nil {
		return err
	}
	return ou.audit.Modification(accountId, orderId)
}
