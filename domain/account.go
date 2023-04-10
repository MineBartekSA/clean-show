package domain

import "unicode"

type Account struct {
	DBModel `json:"-"`
	Type    AccountType `db:"type" json:"type" patch:"-"`
	Email   string      `db:"email" json:"email"`
	Hash    string      `db:"hash" json:"-"`
	Name    string      `db:"name" json:"name"`
	Surname string      `db:"surname" json:"surname"`
}

type AccountType int

const (
	AccountTypeUser AccountType = iota + 1
	AccountTypeStaff
)

type AccountLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccountCreate struct {
	*AccountLogin
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func (al *AccountLogin) Validate() error {
	return ValidatePassword(al.Password)
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrBadPassword.Call()
	}
	up := false
	low := false
	digit := false
	other := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			up = true
		} else if unicode.IsLower(char) {
			low = true
		} else if unicode.IsDigit(char) {
			digit = true
		} else {
			other = true
		}
		if up && low && digit && other {
			break
		}
	}
	if !up || !low || !digit || !other {
		return ErrBadPassword.Call()
	}
	return nil
}

//go:generate mockery --name AccountController
type AccountController interface {
	Register(router Router)
	PostLogin(context Context, session UserSession)
	PostRegister(context Context, session UserSession)
	GetByID(context Context, session UserSession)
	Patch(context Context, session UserSession)
	GetOrders(context Context, session UserSession)
	PostPassword(context Context, session UserSession)
	GetLogout(context Context, session UserSession)
	Delete(context Context, session UserSession)
}

//go:generate mockery --name AccountUsecase
type AccountUsecase interface {
	Login(login *AccountLogin) (*Account, string, error)
	Register(register *AccountCreate) (*Account, string, error)
	FetchBySession(session *Session) (*Account, error)
	FetchByID(session UserSession, id uint) (*Account, error)
	Modify(session UserSession, accountId uint, data map[string]any) (*Account, error)
	FetchOrders(session UserSession, accountId uint, limit, page int) ([]Order, error)
	ModifyPassword(session UserSession, accountId uint, new string) error
	Logout(session UserSession) error
	Remove(session UserSession, accountId uint) error
}

//go:generate mockery --name AccountRepository
type AccountRepository interface {
	SelectEMail(email string) (*Account, error)
	SelectID(id uint, full bool) (*Account, error)
	Insert(account *Account) error
	Update(account *Account) error
	UpdateHash(accountId uint, hash string) error
	Delete(id uint) error
}
