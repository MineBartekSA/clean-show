package account

import (
	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/domain"
)

type accountRepository struct {
	db domain.DB

	selectEMail       domain.Stmt
	selectIDPartial   domain.Stmt
	selectIDFull      domain.Stmt
	insert            domain.Stmt
	update            domain.Stmt
	updateDeleteEmail domain.Stmt
	delete            domain.Stmt
}

func NewAccountRepository(db domain.DB) domain.AccountRepository {
	return &accountRepository{
		db:                db,
		selectEMail:       db.PrepareSelect("accounts", "email = :email"),
		selectIDPartial:   db.Prepare("SELECT type, email, name, surname FROM accounts WHERE deleted_at IS NULL AND id = :id"),
		selectIDFull:      db.PrepareSelect("accounts", "id = :id"),
		insert:            db.PrepareInsertStruct("accounts", &domain.Account{}),
		update:            db.PrepareUpdateStruct("accounts", &domain.Account{}, "id = :id"),
		updateDeleteEmail: db.PrepareUpdate("accounts", "email = '@'+email", "id = :id"),
		delete:            db.PrepareSoftDelete("accounts", "id = :id"),
	}
}

func (ar *accountRepository) SelectEMail(email string) (*domain.Account, error) {
	var account domain.Account
	err := ar.selectEMail.Get(&account, domain.H{"email": email})
	return &account, domain.SQLError(err)
}

func (ar *accountRepository) SelectID(id uint, full bool) (res *domain.Account, err error) {
	var account domain.Account
	if full {
		err = ar.selectIDFull.Get(&account, domain.H{"id": id})
	} else {
		err = ar.selectIDPartial.Get(&account, domain.H{"id": id})
	}
	return &account, domain.SQLError(err)
}

func (ar *accountRepository) Insert(account *domain.Account) error {
	var err error
	if config.Env.DBDriver == "mysql" { // TODO: Try to generalize Inserts
		err = ar.db.Transaction(func(tx domain.Tx) error {
			res, err := tx.Stmt(ar.insert).Exec(account)
			if err != nil {
				return err
			}
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			account.ID = uint(id)
			return nil
		})
	} else {
		err = ar.insert.Get(account, account)
	}
	return domain.SQLError(err)
}

func (ar *accountRepository) Update(account *domain.Account) error { // TODO: Hash specific Update method?
	_, err := ar.update.Exec(account)
	return domain.SQLError(err)
}

func (ar *accountRepository) Delete(id uint) error {
	return domain.SQLError(ar.db.Transaction(func(tx domain.Tx) error {
		_, err := tx.Stmt(ar.updateDeleteEmail).Exec(domain.H{"id": id})
		if err != nil {
			return err
		}
		_, err = tx.Stmt(ar.delete).Exec(domain.H{"id": id})
		return err
	}))
}
