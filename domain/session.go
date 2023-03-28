package domain

type Session struct {
	DBModel
	AccountID uint   `db:"account_id"`
	Token     string `db:"token"`
}

type UserSession interface {
	// TODO: fill
}

type SessionUsecase interface {
	Fetch(token string) (*Session, error)
	Create(account_id uint) (*Session, error)
	Invalidate(session_id uint) error
}

type SessionRepository interface {
	SelectByToken(token string) (*Session, error)
	Insert(session Session) error
	Delete(session_id uint) error
}
