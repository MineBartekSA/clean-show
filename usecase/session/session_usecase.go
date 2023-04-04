package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"time"

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
	buffer := make([]byte, 128)
	_, err := rand.Read(buffer)
	if err != nil {
		return nil, err
	}
	randString := base64.URLEncoding.EncodeToString(buffer)
	buffer = make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buffer, time.Now().UTC().UnixMilli())
	token := base64.URLEncoding.EncodeToString(buffer[:n]) + "."
	token += randString[:128-len(token)]
	session := domain.Session{
		AccountID: account_id,
		Token:     token,
	}
	err = su.repository.Insert(&session)
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
	return su.audit.Deletion(session.AccountID, session.ID)
}

func (su *sessionUsecase) InvalidateAccount(executorId, accountId uint) error {
	err := su.repository.DeleteByAccount(accountId)
	if err != nil {
		return err
	}
	return su.audit.Deletion(executorId, 0)
}
