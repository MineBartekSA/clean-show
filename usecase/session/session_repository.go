package session

import (
	"time"

	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
)

type sessionRespository struct {
	db domain.DB

	tokenSelect domain.Stmt
	insert      domain.Stmt
	update      domain.Stmt
	delete      domain.Stmt
}

func NewSessionRepository(db domain.DB) domain.SessionRepository {
	tokenSelect, err := db.PrepareSelect("sessions", "token = :token AND updated_at > "+domain.DBInterval(domain.DBNow(), time.Minute*-30))
	if err != nil {
		Log.Panicw("failed to prepare a named select statement", "err", err)
	}
	insert, err := db.PrepareInsertStruct("sessions", &domain.Session{})
	if err != nil {
		Log.Panicw("failed to prepare a named insert statement from a structure", "err", err)
	}
	update, err := db.PrepareUpdate("sessions", "", "id = :id")
	if err != nil {
		Log.Panicw("failed to prepare a named update statement", "err", err)
	}
	delete, err := db.PrepareSoftDelete("sessions", "id = :id")
	if err != nil {
		Log.Panicw("failed to prepare a named soft delete statement", "err", err)
	}
	return &sessionRespository{db, tokenSelect, insert, update, delete}
}

func (sr *sessionRespository) SelectByToken(token string) (*domain.Session, error) {
	var session domain.Session
	err := sr.tokenSelect.Get(&session, domain.H{"token": token})
	return &session, err
}

func (sr *sessionRespository) Insert(session *domain.Session) error {
	var err error
	if config.Env.DBDriver == "mysql" {
		err = sr.db.Transaction(func(tx domain.Tx) error {
			res, err := tx.Stmt(sr.insert).Exec(sr.db.PrepareStruct(&session))
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
		err = sr.insert.Get(&session, sr.db.PrepareStruct(&session))
	}
	return err
}

func (sr *sessionRespository) Extend(session_id uint) error {
	_, err := sr.update.Exec(domain.H{"id": session_id})
	return err
}

func (sr *sessionRespository) Delete(session_id uint) error {
	_, err := sr.delete.Exec(domain.H{"id": session_id})
	return err
}
