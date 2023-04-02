package domain

type Account struct {
	DBModel `json:"-"`
	Type    AccountType `db:"type" json:"type" patch:"-"`
	Email   string      `db:"email" json:"email"`
	Hash    string      `db:"hash" json:"-"`
	Salt    string      `db:"salt" json:"-"`
	Name    string      `db:"name" json:"name"`
	Surname string      `db:"surname" json:"surname"`
}

type AccountType int

const (
	AccountTypeUser AccountType = iota + 1
	AccountTypeStaff
)

type AccountController interface {
	Register(router Router)
	GetByID(context Context, session UserSession)
}

type AccountUsecase interface {
	FetchBySession(session *Session) (*Account, error)
	FetchByID(session UserSession, id uint) (*Account, error)
}

type AccountRepository interface {
	SelectID(id uint, full bool) (*Account, error)
}
