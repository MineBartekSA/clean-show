package product

import (
	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
)

type productRepository struct {
	db domain.DB

	count      domain.Stmt
	selectList domain.Stmt
	selectID   domain.Stmt
	insert     domain.Stmt
	update     domain.Stmt
	delete     domain.Stmt
}

func NewProductRepository(db domain.DB) domain.ProductRepository {
	count, err := db.Prepare("SELECT COUNT(*) FROM products")
	if err != nil {
		Log.Panicw("failed to prepare a named count select statement", "err", err)
	}
	selectList, err := db.Prepare("SELECT * FROM products WHERE deleted_at IS NULL LIMIT :limit OFFSET :offset")
	if err != nil {
		Log.Panicw("failed to prepare a named select statement", "err", err)
	}
	selectID, err := db.PrepareSelect("products", "id = :id")
	if err != nil {
		Log.Panicw("failed to prepare a named select statement", "err", err)
	}
	insert, err := db.PrepareInsertStruct("products", &domain.Product{})
	if err != nil {
		Log.Panicw("failed to prepare a named insert statement from structure", "err", err)
	}
	update, err := db.PrepareUpdateStruct("products", &domain.Product{}, "id = :id")
	if err != nil {
		Log.Panicw("failed to prepare a named update statement from structure", "err", err)
	}
	delete, err := db.PrepareSoftDelete("products", "id = :id")
	if err != nil {
		Log.Panicw("failed to prepare a named soft delete statement", "err", err)
	}
	return &productRepository{db, count, selectList, selectID, insert, update, delete}
}

func (pr *productRepository) Count() (res uint, err error) {
	err = pr.count.Get(&res, domain.H{})
	return
}

func (pr *productRepository) Select(limit, page int) ([]domain.Product, error) {
	res := []domain.Product{}
	err := pr.selectList.Select(&res, domain.H{"limit": limit, "offset": (page - 1) * limit})
	return res, err
}

func (pr *productRepository) SelectID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := pr.selectID.Get(&product, domain.H{"id": id})
	return &product, err
}

func (pr *productRepository) Insert(product *domain.Product) error {
	var err error
	if config.Env.DBDriver == "mysql" {
		err = pr.db.Transaction(func(tx domain.Tx) error {
			res, err := tx.Stmt(pr.insert).Exec(pr.db.PrepareStruct(&product))
			if err != nil {
				return err
			}
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			product.ID = uint(id)
			return nil
		})
	} else {
		err = pr.insert.Get(&product, pr.db.PrepareStruct(&product))
	}
	return err
}

func (pr *productRepository) Update(product *domain.Product) error {
	_, err := pr.update.Exec(product)
	return err
}

func (pr *productRepository) Delete(id uint) error {
	_, err := pr.delete.Exec(domain.H{"id": id})
	return err
}
