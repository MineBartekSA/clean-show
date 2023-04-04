package product

import (
	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/domain"
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
	return &productRepository{
		db:         db,
		count:      db.Prepare("SELECT COUNT(*) FROM products"),
		selectList: db.Prepare("SELECT * FROM products WHERE deleted_at IS NULL LIMIT :limit OFFSET :offset"),
		selectID:   db.PrepareSelect("products", "id = :id"),
		insert:     db.PrepareInsertStruct("products", &domain.Product{}),
		update:     db.PrepareUpdateStruct("products", &domain.Product{}, "id = :id"),
		delete:     db.PrepareSoftDelete("products", "id = :id"),
	}
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
	if config.Env.DBDriver == "mysql" { // TODO: Try to generalize Inserts
		err = pr.db.Transaction(func(tx domain.Tx) error {
			res, err := tx.Stmt(pr.insert).Exec(pr.db.PrepareStruct(product))
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
		err = pr.insert.Get(product, pr.db.PrepareStruct(product))
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
