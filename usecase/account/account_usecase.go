package account

import "github.com/minebarteksa/clean-show/domain"

type accountUsecase struct {
	repository    domain.AccountRepository
	auditPassword domain.AuditResource
	audit         domain.AuditResource
}

func NewAccountUsecase(repository domain.AccountRepository, audit domain.AuditUsecase) domain.AccountUsecase {
	return &accountUsecase{repository, audit.Resource(domain.ResourceTypeAccountPassword), audit.Resource(domain.ResourceTypeAccount)}
}

func (au *accountUsecase) FetchBySession(session *domain.Session) (*domain.Account, error) {
	return au.repository.SelectID(session.AccountID)
}
