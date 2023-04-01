package product

import (
	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
)

type productRepository struct {
	db domain.DB

	selectID domain.Stmt
	insert   domain.Stmt
}

func NewProductRepository(db domain.DB) domain.ProductRepository {
	selectID, err := db.PrepareSelect("products", "id = :id")
	if err != nil {
		Log.Panicw("failed to prepare a named select statement", "err", err)
	}
	insert, err := db.PrepareInsertStruct("products", &domain.Product{})
	if err != nil {
		Log.Panicw("failed to prepare a named insert statement from structure", "err", err)
	}
	return &productRepository{db, selectID, insert}
}

func (pr *productRepository) ID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := pr.selectID.Select(&product, &domain.H{"id": id})
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
		err = pr.insert.Select(&product, pr.db.PrepareStruct(&product))
	}
	return err
}
