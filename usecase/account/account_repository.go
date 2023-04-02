package account

import (
	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
)

type accountRepository struct {
	db domain.DB

	selectPartial domain.Stmt
	selectFull    domain.Stmt
}

func NewAccountRepository(db domain.DB) domain.AccountRepository {
	selectPartial, err := db.Prepare("SELECT type, email, name, surname FROM accounts WHERE deleted_at IS NULL AND id = :id")
	if err != nil {
		Log.Panicw("failed to prepare a named select statement", "err", err)
	}
	selectFull, err := db.PrepareSelect("accounts", "id = :id")
	if err != nil {
		Log.Panicw("failed to prepare a named select statement", "err", err)
	}
	return &accountRepository{db, selectPartial, selectFull}
}

func (ar *accountRepository) SelectID(id uint, full bool) (res *domain.Account, err error) {
	var account domain.Account
	if full {
		err = ar.selectFull.Get(&account, domain.H{"id": id})
	} else {
		err = ar.selectPartial.Get(&account, domain.H{"id": id})
	}
	return &account, err
}
