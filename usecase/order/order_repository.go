package order

import (
	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/domain"
)

type orderRepository struct {
	db domain.DB

	count         domain.Stmt
	selectList    domain.Stmt
	selectAccount domain.Stmt
	selectID      domain.Stmt
	selectOrderBy domain.Stmt
	insert        domain.Stmt
	update        domain.Stmt
	updateStatus  domain.Stmt
	delete        domain.Stmt
}

func NewOrderRepository(db domain.DB) domain.OrderRepository {
	return &orderRepository{
		db:            db,
		count:         db.Prepare("SELECT COUNT(*) FROM orders"),
		selectList:    db.Prepare("SELECT * FROM orders WHERE deleted_at IS NULL ORDER BY id DESC LIMIT :limit OFFSET :offset"),
		selectAccount: db.PrepareSelect("orders", "order_by = :account"),
		selectID:      db.PrepareSelect("orders", "id = :id"),
		selectOrderBy: db.Prepare("SELECT order_by FROM orders WHERE id = :id AND deleted_at IS NULL"),
		insert:        db.PrepareInsertStruct("orders", &domain.Order{}),
		update:        db.PrepareUpdateStruct("orders", &domain.Order{}, "id = :id"),
		updateStatus:  db.PrepareUpdate("orders", "status = :status", "id = :id"),
		delete:        db.PrepareSoftDelete("orders", "id = :id"),
	}
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

func (or *orderRepository) SelectAccount(accountId uint) ([]domain.Order, error) {
	res := []domain.Order{}
	err := or.selectAccount.Select(&res, domain.H{"account": accountId})
	return res, err
}

func (or *orderRepository) SelectID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := or.selectID.Get(&order, domain.H{"id": id})
	return &order, err
}

func (or *orderRepository) SelectOrderBy(orderId uint) (uint, error) {
	var by uint
	err := or.selectOrderBy.Get(&by, domain.H{"id": orderId})
	return by, err
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

func (or *orderRepository) UpdateStatus(orderId uint, status domain.OrderStatus) error {
	_, err := or.updateStatus.Exec(domain.H{"id": orderId, "status": status})
	return err
}

func (or *orderRepository) Delete(id uint) error {
	_, err := or.delete.Exec(domain.H{"id": id})
	return err
}
