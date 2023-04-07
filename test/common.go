package test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/domain/mocks"
	"github.com/minebarteksa/clean-show/infrastructure/database"
	"github.com/stretchr/testify/assert"
)

func SetupConfig() {
	config.Env = &config.EnvConfig{
		Debug:    true,
		Port:     8080,
		DBDriver: "postgres",
		DBSource: "none lmao",
	}
}

func NewMockDB(t *testing.T) (domain.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS audit_log\\(.*\\)").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("CREATE INDEX IF NOT EXISTS .*").WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS sessions\\(.*\\)").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("CREATE INDEX IF NOT EXISTS .*").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("CREATE INDEX IF NOT EXISTS .*").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("CREATE UNIQUE INDEX IF NOT EXISTS .*").WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS accounts\\(.*\\)").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("CREATE INDEX IF NOT EXISTS .*").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("CREATE UNIQUE INDEX IF NOT EXISTS .*").WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS products\\(.*\\)").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("CREATE INDEX IF NOT EXISTS .*").WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS orders\\(.*\\)").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("CREATE INDEX IF NOT EXISTS .*").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("CREATE INDEX IF NOT EXISTS .*").WillReturnResult(sqlmock.NewResult(0, 1))

	return database.NewSqlDB(sqlx.NewDb(db, "sqlmock")), mock
}

func NewRows(columns ...string) *sqlmock.Rows {
	return sqlmock.NewRows(columns)
}

func IDRow(id uint) *sqlmock.Rows {
	return NewRows("id").AddRow(id)
}

func NewAuditUsecase(t *testing.T, resource domain.ResourceType) (*mocks.AuditResource, domain.AuditUsecase) {
	mock := mocks.NewAuditUsecase(t)
	res := mocks.NewAuditResource(t)
	mock.On("Resource", resource).Return(res)
	return res, mock
}
