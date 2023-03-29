package router

import (
	"github.com/minebarteksa/clean-show/domain"
)

type userSession struct {
	session *domain.Session
	account *domain.Account
}

func EmptySession() domain.UserSession {
	return &userSession{}
}

func NewSession(session *domain.Session, account *domain.Account) domain.UserSession {
	return &userSession{session, account}
}

func (us *userSession) Authorized() bool {
	return us.session != nil
}

func (us *userSession) GetAccount() *domain.Account {
	return us.account
}
