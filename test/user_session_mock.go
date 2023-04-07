package test

import "github.com/minebarteksa/clean-show/domain"

type MockUserSession struct {
	*domain.Session
	*domain.Account
	Authed bool
}

func (usm *MockUserSession) Authorized() bool {
	return usm.Authed
}

func (usm *MockUserSession) GetSession() *domain.Session {
	return usm.Session
}

func (usm *MockUserSession) GetAccount() *domain.Account {
	return usm.Account
}

func (usm *MockUserSession) GetAccountID() uint {
	return usm.Account.ID
}

func (usm *MockUserSession) IsStaff() bool {
	return usm.Account.Type == domain.AccountTypeStaff
}
