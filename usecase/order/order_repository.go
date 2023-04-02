package order

import (
	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
)

type orderRepository struct {
	db domain.DB

	count      domain.Stmt
	selectList domain.Stmt
	selectID   domain.Stmt
	insert     domain.Stmt
	update     domain.Stmt
	delete     domain.Stmt
}

func NewOrderRepository(db domain.DB) domain.OrderRepository {
	count, err := db.Prepare("SELECT COUNT(*) FROM orders")
	if err != nil {
		Log.Panicw("failed to prepare a named count select statement", "err", err)
	}
	selectList, err := db.Prepare("SELECT * FROM orders WHERE deleted_at IS NULL ORDER BY id DESC LIMIT :limit OFFSET :offset")
	if err != nil {
		Log.Panicw("failed to prepare a named select statement", "err", err)
	}
	selectID, err := db.PrepareSelect("orders", "id = :id")
	if err != nil {
		Log.Panicw("failed to prepare a named select statement", "err", err)
	}
	insert, err := db.PrepareInsertStruct("orders", &domain.Order{})
	if err != nil {
		Log.Panicw("failed to prepare a named insert statement from structure", "err", err)
	}
	update, err := db.PrepareUpdateStruct("orders", &domain.Order{}, "id = :id")
	if err != nil {
		Log.Panicw("failed to prepare a named update statement from structure", "err", err)
	}
	delete, err := db.PrepareSoftDelete("orders", "id = :id")
	if err != nil {
		Log.Panicw("failed to prepare a named soft delete statement", "err", err)
	}
	return &orderRepository{db, count, selectList, selectID, insert, update, delete}
}

func (or *orderRepository) Count() (res uint, err error) {
	err = or.count.Get(&res, domain.H{})
	return
}

func (or *orderRepository) Select(limit, page int) ([]domain.Order, error) {
	res := []domain.Order{}
	err := or.selectList.Select(&res, domain.H{"limit": limit, "offset": (page - 1) * limit})
	return res, err
}

func (or *orderRepository) SelectID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := or.selectID.Get(&order, domain.H{"id": id})
	return &order, err
}

func (or *orderRepository) Insert(order *domain.Order) error {
	var err error
	if config.Env.DBDriver == "mysql" { // TODO: Try to generalize Inserts
		err = or.db.Transaction(func(tx domain.Tx) error {
			res, err := tx.Stmt(or.insert).Exec(or.db.PrepareStruct(order))
			if err != nil {
				return err
			}
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			order.ID = uint(id)
			return nil
		})
	} else {
		err = or.insert.Get(order, or.db.PrepareStruct(order))
	}
	return err
}

func (or *orderRepository) Update(order *domain.Order) error {
	_, err := or.update.Exec(order)
	return err
}

func (or *orderRepository) Delete(id uint) error {
	_, err := or.delete.Exec(domain.H{"id": id})
	return err
}
