package domain

import "database/sql"

type Stmt interface {
	Exec(args ...any) (sql.Result, error)
	Select(dst interface{}, args ...interface{}) error
}

type Tx interface {
	Stmt(Stmt) Stmt
}

type DB interface {
	Prepare(query string) (Stmt, error)
	Exec(query string, args ...any) (sql.Result, error)
	Select(dst interface{}, query string, args ...interface{}) error

	Transaction(func(tx Tx) error) error
}
