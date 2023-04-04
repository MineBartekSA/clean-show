package session

import (
	"time"

	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/domain"
)

type sessionRespository struct {
	db domain.DB

	tokenSelect   domain.Stmt
	insert        domain.Stmt
	update        domain.Stmt
	delete        domain.Stmt
	deleteAccount domain.Stmt
}

func NewSessionRepository(db domain.DB) domain.SessionRepository {
	return &sessionRespository{
		db:            db,
		tokenSelect:   db.PrepareSelect("sessions", "updated_at > "+domain.DBInterval(domain.DBNow(), time.Minute*-30)+" AND token = :token"),
		insert:        db.PrepareInsertStruct("sessions", &domain.Session{}),
		update:        db.PrepareUpdate("sessions", "", "id = :id"),
		delete:        db.PrepareSoftDelete("sessions", "id = :id"),
		deleteAccount: db.PrepareSoftDelete("sessions", "account_id = :account"),
	}
}

func (sr *sessionRespository) SelectByToken(token string) (*domain.Session, error) {
	var session domain.Session
	err := sr.tokenSelect.Get(&session, domain.H{"token": token})
	return &session, domain.SQLError(err)
}

func (sr *sessionRespository) Insert(session *domain.Session) error {
	var err error
	if config.Env.DBDriver == "mysql" {
		err = sr.db.Transaction(func(tx domain.Tx) error {
			res, err := tx.Stmt(sr.insert).Exec(session)
			if err != nil {
				return err
			}
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			session.ID = uint(id)
			return nil
		})
	} else {
		err = sr.insert.Get(&session, session)
	}
	return domain.SQLError(err)
}

func (sr *sessionRespository) Extend(sessionId uint) error {
	_, err := sr.update.Exec(domain.H{"id": sessionId})
	return domain.SQLError(err)
}

func (sr *sessionRespository) Delete(sessionId uint) error {
	_, err := sr.delete.Exec(domain.H{"id": sessionId})
	return domain.SQLError(err)
}

func (sr *sessionRespository) DeleteByAccount(acountId uint) error {
	_, err := sr.deleteAccount.Exec(domain.H{"account": acountId})
	return domain.SQLError(err)
}
