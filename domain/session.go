package domain

type Session struct {
	DBModel
	AccountID uint   `db:"account_id"`
	Token     string `db:"token"`
}

type UserSession interface {
	Authorized() bool
	GetAccount() *Account
}

type SessionUsecase interface {
	Fetch(token string) (*Session, error)
	Create(account_id uint) (*Session, error)
	Invalidate(session *Session) error
}

type SessionRepository interface {
	SelectByToken(token string) (*Session, error)
	Insert(session *Session) error
	Extend(session_id uint) error
	Delete(session_id uint) error
}
