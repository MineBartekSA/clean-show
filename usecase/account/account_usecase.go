package account

import (
	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/usecase"
)

type accountUsecase struct {
	repository domain.AccountRepository

	order   domain.OrderUsecase
	session domain.SessionUsecase

	audit         domain.AuditResource
	auditPassword domain.AuditResource

	hasher domain.Hasher
}

func NewAccountUsecase(repository domain.AccountRepository, order domain.OrderUsecase, session domain.SessionUsecase, audit domain.AuditUsecase, hasher domain.Hasher) domain.AccountUsecase {
	return &accountUsecase{
		repository:    repository,
		order:         order,
		session:       session,
		audit:         audit.Resource(domain.ResourceTypeAccount),
		auditPassword: audit.Resource(domain.ResourceTypeAccountPassword),
		hasher:        hasher,
	}
}

func (au *accountUsecase) Login(login *domain.AccountLogin) (*domain.Account, string, error) {
	account, err := au.repository.SelectEMail(login.Email)
	if err != nil {
		return nil, "", err
	}

	verified, err := au.hasher.Verify(login.Password, account.Hash)
	if err != nil {
		return nil, "", err
	} else if !verified {
		return nil, "", domain.Fatal(domain.ErrUnauthorized, "hash verification failed").Call()
	}

	session, err := au.session.Create(account.ID)
	if err != nil {
		return nil, "", err
	}
	return account, session.Token, nil
}

func (au *accountUsecase) Register(register *domain.AccountCreate) (*domain.Account, string, error) {
	account := domain.Account{
		Type:    domain.AccountTypeUser,
		Email:   register.Email,
		Name:    register.Name,
		Surname: register.Surname,
	}

	err := register.Validate()
	if err != nil {
		return nil, "", err
	}

	account.Hash = au.hasher.Hash(register.Password)

	err = au.repository.Insert(&account)
	if err != nil {
		return nil, "", err
	}
	session, err := au.session.Create(account.ID)
	if err != nil {
		return nil, "", err
	}
	return &account, session.Token, nil
}

func (au *accountUsecase) FetchBySession(session *domain.Session) (*domain.Account, error) {
	return au.repository.SelectID(session.AccountID, true)
}

func (au *accountUsecase) FetchByID(session domain.UserSession, id uint) (*domain.Account, error) {
	if !session.IsStaff() && id != session.GetAccountID() {
		return nil, domain.Fatal(domain.ErrUnauthorized, "only staff users can fetch other accounts information").Call()
	}
	return au.repository.SelectID(id, false)
}

func (au *accountUsecase) Modify(session domain.UserSession, accountId uint, data map[string]any) (account *domain.Account, err error) {
	account = session.GetAccount()
	if !session.IsStaff() && account.ID != accountId {
		return nil, domain.Fatal(domain.ErrUnauthorized, "only staff users can modify other accounts information").Call()
	}
	if account.ID != accountId {
		account, err = au.repository.SelectID(accountId, false)
		if err != nil {
			return nil, err
		}
		if account.Type == domain.AccountTypeStaff {
			return nil, domain.Fatal(domain.ErrUnauthorized, "only the owner of this account can change its information").Call()
		}
	}
	err = usecase.PatchModel(account, data)
	if err != nil {
		return nil, err
	}
	err = au.repository.Update(account)
	if err != nil {
		return nil, err
	}
	return account, au.audit.Modification(session.GetAccountID(), accountId)
}

func (au *accountUsecase) FetchOrders(session domain.UserSession, accountId uint, limit, page int) ([]domain.Order, error) {
	if !session.IsStaff() && session.GetAccountID() != accountId {
		return nil, domain.Fatal(domain.ErrUnauthorized, "only staff users can fetch other users orders").Call()
	}
	return au.order.FetchByAccount(accountId, limit, page)
}

func (au *accountUsecase) ModifyPassword(session domain.UserSession, accountId uint, new string) error {
	aid := session.GetAccountID()
	if !session.IsStaff() && aid != accountId {
		return domain.Fatal(domain.ErrUnauthorized, "only staff users can modify other users passwords").Call()
	}
	account, err := au.repository.SelectID(accountId, false)
	if err != nil {
		return err
	}
	if aid != account.ID && account.Type == domain.AccountTypeStaff {
		return domain.Fatal(domain.ErrUnauthorized, "only the owner of this account can change its password").Call()
	}

	err = domain.ValidatePassword(new)
	if err != nil {
		return err
	}
	err = au.repository.UpdateHash(account.ID, au.hasher.Hash(new))
	if err != nil {
		return err
	}
	return au.auditPassword.Modification(aid, accountId)
}

func (au *accountUsecase) Logout(session domain.UserSession) error {
	return au.session.Invalidate(session.GetSession())
}

func (au *accountUsecase) Remove(session domain.UserSession, accountId uint) error {
	aid := session.GetAccountID()
	if !session.IsStaff() && aid != accountId {
		return domain.Fatal(domain.ErrUnauthorized, "only staff users can remove other users accounts").Call()
	}
	err := au.repository.Delete(accountId)
	if err != nil {
		return err
	}
	err = au.order.CancelByAccount(aid, accountId)
	if err != nil {
		return err
	}
	err = au.session.InvalidateAccount(aid, accountId)
	if err != nil {
		return err
	}
	return au.audit.Deletion(aid, accountId)
}
