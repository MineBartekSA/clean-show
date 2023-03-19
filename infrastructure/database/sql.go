package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
)

type sqldb struct {
	*sqlx.DB
}

func NewSqlDB() domain.DB {
	db, err := sqlx.Connect(config.Env.DBDriver, config.Env.DBDriver)
	if err != nil {
		Log.Fatalw("failed to connect to the database", "err", err)
	}
	return &sqldb{db}
}

func (db *sqldb) Prepare(query string) (domain.Stmt, error) {
	if config.Env.DBDriver == "postgres" || config.Env.DBDriver == "sqlite3" {
		index := 1
		for i := 0; i < len(query); i += 1 {
			if query[i] == '?' {
				placeholder := fmt.Sprintf("$%d", index)
				query = query[:i] + placeholder + query[i+len(placeholder):]
				i += len(placeholder)
				index += 1
			}
		}
	}
	return db.Preparex(query)
}

func (db *sqldb) Transaction(transaction func(domain.Tx) error) error {
	tx, err := db.DB.Beginx()
	err = transaction(&sqltx{tx})
	if err != nil {
		e := tx.Rollback()
		if e != nil {
			return e
		}
		return err
	}
	return tx.Commit()
}

type sqltx struct {
	*sqlx.Tx
}

func (tx *sqltx) Stmt(stmt domain.Stmt) domain.Stmt {
	return tx.Tx.Stmtx(stmt)
}
