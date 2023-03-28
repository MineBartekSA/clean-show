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
	AccountTypeStaff AccountType = iota + 1
	AccountTypeUser
)

// TODO: Implement interfaces
