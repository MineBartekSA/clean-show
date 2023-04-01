package account

import (
	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
)

type accountRepository struct {
	db domain.DB

	selectID domain.Stmt
}

func NewAccountRepository(db domain.DB) domain.AccountRepository {
	selectID, err := db.PrepareSelect("accounts", "id = :id")
	if err != nil {
		Log.Panicw("failed to prepare a named select statement", "err", err)
	}
	return &accountRepository{db, selectID}
}

func (ar *accountRepository) SelectID(id uint) (*domain.Account, error) {
	var account domain.Account
	err := ar.selectID.Get(&account, domain.H{"id": id})
	return &account, err
}
