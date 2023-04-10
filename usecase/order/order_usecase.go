package order

import (
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

func (ou *orderUsecase) FetchByAccount(accountId uint, limit, page int) ([]domain.Order, error) {
	if limit < 0 {
		limit = 0
	} else if limit > 1000 {
		limit = 1000
	}
	if page < 1 {
		page = 1
	}
	return ou.repository.SelectAccount(accountId, limit, page)
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
		return nil, domain.Fatal(domain.ErrUnauthorized, "only staff users can see other users orders").Call()
	}
	return order, nil
}

func (ou *orderUsecase) Modify(accountId, orderId uint, data map[string]any) (*domain.Order, error) {
	order, err := ou.repository.SelectID(orderId)
	if err != nil {
		return nil, err
	}
	err = usecase.PatchModel(order, data)
	if err != nil {
		return nil, err
	}
	order.UpdateTotal()
	err = ou.repository.Update(order)
	if err != nil {
		return nil, err
	}
	return order, ou.audit.Modification(accountId, orderId)
}

func (ou *orderUsecase) Cancel(session domain.UserSession, orderId uint) error {
	aid := session.GetAccountID()
	orderBy, err := ou.repository.SelectOrderBy(orderId)
	if err != nil {
		return err
	}
	if !session.IsStaff() && aid != orderBy {
		return domain.Fatal(domain.ErrUnauthorized, "only staff users can cancel other users orders").Call()
	}
	err = ou.repository.UpdateStatus(orderId, domain.OrderStatusCanceled)
	if err != nil {
		return err
	}
	return ou.audit.Modification(aid, orderId)
}

func (ou *orderUsecase) CancelByAccount(executorId, accountId uint) error {
	orders, err := ou.repository.SelectAccount(accountId, 0, 0)
	if err != nil {
		return err
	}
	var orderIds []uint
	for _, order := range orders {
		if order.Status == domain.OrderStatusCompleted || order.Status == domain.OrderStatusCanceled {
			continue
		}
		orderIds = append(orderIds, order.ID)
	}
	err = ou.repository.BatchUpdateStatus(orderIds, domain.OrderStatusCanceled)
	if err != nil {
		return err
	}
	err = ou.audit.BatchModification(executorId, orderIds)
	if err != nil {
		return err
	}
	return nil
}

func (ou *orderUsecase) Remove(accountId, orderId uint) error {
	err := ou.repository.Delete(orderId)
	if err != nil {
		return err
	}
	return ou.audit.Deletion(accountId, orderId)
}
