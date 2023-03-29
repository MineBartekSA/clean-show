package session

import (
	"github.com/dchest/uniuri"
	"github.com/minebarteksa/clean-show/domain"
)

type sessionUsecase struct {
	repository domain.SessionRepository
	audit      domain.AuditResource
}

func NewSessionUsecase(repository domain.SessionRepository, audit domain.AuditUsecase) domain.SessionUsecase {
	return &sessionUsecase{repository, audit.Resource(domain.ResourceTypeSession)}
}

func (su *sessionUsecase) Fetch(token string) (*domain.Session, error) {
	session, err := su.repository.SelectByToken(token)
	if err != nil {
		return nil, err
	}
	err = su.repository.Extend(session.ID)
	if err != nil {
		return nil, err
	}
	err = su.audit.Modification(session.AccountID, session.ID)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (su *sessionUsecase) Create(account_id uint) (*domain.Session, error) {
	session := domain.Session{
		AccountID: account_id,
		Token:     uniuri.NewLen(128),
	}
	err := su.repository.Insert(&session)
	if err != nil {
		return nil, err
	}
	err = su.audit.Creation(session.AccountID, session.ID)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (su *sessionUsecase) Invalidate(session *domain.Session) error {
	err := su.repository.Delete(session.ID)
	if err != nil {
		return err
	}
	err = su.audit.Deletion(session.AccountID, session.ID)
	if err != nil {
		return err
	}
	return nil
}
