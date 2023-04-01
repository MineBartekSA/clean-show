package domain

type Account struct {
	DBModel
	Type    AccountType `db:"type"`
	Email   string      `db:"email"`
	Hash    string      `db:"hash"`
	Salt    string      `db:"salt"`
	Name    string      `db:"name"`
	Surname string      `db:"surname"`
}

type AccountType int

const (
	AccountTypeUser AccountType = iota + 1
	AccountTypeStaff
)

type AccountController interface {
	// TODO: Implement
}

type AccountUsecase interface {
	FetchBySession(session *Session) (*Account, error)
}

type AccountRepository interface {
	SelectID(id uint) (*Account, error)
}
