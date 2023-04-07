package domain

type Session struct {
	DBModel
	AccountID uint   `db:"account_id"`
	Token     string `db:"token"`
}

type UserSession interface {
	Authorized() bool
	GetSession() *Session
	GetAccount() *Account
	GetAccountID() uint
	IsStaff() bool
}

//go:generate mockery --name SessionUsecase
type SessionUsecase interface {
	Fetch(token string) (*Session, error)
	Create(account_id uint) (*Session, error)
	Invalidate(session *Session) error
	InvalidateAccount(executorId, accountId uint) error
}

//go:generate mockery --name SessionRepository
type SessionRepository interface {
	SelectByToken(token string) (*Session, error)
	Insert(session *Session) error
	Extend(sessionId uint) error
	Delete(sessionId uint) error
	DeleteByAccount(acountId uint) error
}
