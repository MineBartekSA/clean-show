package database

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
	db, err := sqlx.Connect(config.Env.DBDriver, config.Env.DBSource)
	if err != nil {
		Log.Fatalw("failed to connect to the database", "err", err)
	}

	err = db.Ping()
	if err != nil {
		Log.Fatalw("failed to ping databse", "err", err)
	}

	db.MustExec(createTable("audit_log", `
		type INT NOT NULL,
		resource_type INT NOT NULL,
		resource_id INT NOT NULL,
		executor_id INT NOT NULL
	`))
	db.MustExec(createTable("sessions", `
		account_id INT NOT NULL,
		token TEXT NOT NULL
	`))
	db.MustExec(createTable("accounts", `
		type INT NOT NULL,
		email TEXT NOT NULL,
		hash TEXT NOT NULL,
		salt TEXT NOT NULL,
		name TEXT NOT NULL,
		surname TEXT NOT NULL
	`))
	db.MustExec(createTable("products", `
		status INT NOT NULL,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		price REAL NOT NULL,
		images TEXT NOT NULL
	`))

	return &sqldb{db}
}

func (db *sqldb) PrepareStruct(arg any) any { // TODO: Maybe do this automaticly // Is this needed? DB does this automatically. UpdatedAt is built in into the query
	val := reflect.ValueOf(arg)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		Log.Panicw("argument not a structure")
	}
	dbm := val.FieldByName("DBModel")
	if dbm == reflect.ValueOf(nil) {
		Log.Panicw("argument structure does not contain DBModel")
	}
	now := reflect.ValueOf(time.Now())
	created := dbm.FieldByName("CreatedAt")
	if created.IsZero() {
		created.Set(now)
	}
	dbm.FieldByName("UpdatedAt").Set(now)
	return arg
}

func (db *sqldb) PrepareInsertStruct(table string, arg any) (domain.Stmt, error) {
	cols := getStructureKeys(arg)
	sql := ""
	if config.Env.DBDriver == "mysql" {
		sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (:%s)", table, strings.Join(cols, ", "), strings.Join(cols, ", :"))
	} else {
		sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (:%s) RETURNING id", table, strings.Join(cols, ", "), strings.Join(cols, ", :"))
	}
	return db.Prepare(sql)
}

func (db *sqldb) PrepareSelect(table string, where string) (domain.Stmt, error) {
	if where != "" {
		where += " AND "
	}
	where += "deleted_at IS NULL"
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, where)
	return db.Prepare(sql)
}

func (db *sqldb) PrepareUpdateStruct(table string, arg any, where string) (domain.Stmt, error) {
	cols := getStructureKeys(arg)
	set := ""
	for _, col := range cols {
		set += fmt.Sprintf("%s = :%s, ", col, col)
	}
	return db.PrepareUpdate(table, set[:len(set)-2], where)
}

func (db *sqldb) PrepareUpdate(table, set, where string) (domain.Stmt, error) {
	if set != "" {
		set += ", "
	}
	set += "updated_at = " + domain.DBNow()
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, set, where)
	return db.Prepare(sql)
}

func (db *sqldb) PrepareSoftDelete(table string, where string) (domain.Stmt, error) {
	sql := fmt.Sprintf("UPDATE %s SET deleted_at = %s WHERE %s", table, domain.DBNow(), where)
	return db.Prepare(sql)
}

func (db *sqldb) Prepare(query string) (domain.Stmt, error) {
	Log.Debugf("Preparing: '%s'", query)
	return db.PrepareNamed(query)
}

func (db *sqldb) Transaction(transaction func(domain.Tx) error) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		Log.Panicw("failed to begin transaction", "err", err)
	}
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
	return tx.Tx.NamedStmt(stmt.(*sqlx.NamedStmt))
}

// Utils

func getStructureKeys(arg any) []string { // TODO: Test
	val := reflect.ValueOf(arg)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	var keys []string
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		key := strings.ToLower(field.Name)
		tag := field.Tag.Get("db")
		if tag != "" {
			key = tag
		}
		if key == "-" || key == "dbmodel" {
			continue
		}
		f := val.Field(i)
		switch f.Kind() {
		case reflect.Slice:
			if !f.Type().Implements(reflect.TypeOf((*driver.Valuer)(nil)).Elem()) {
				continue
			}
		case reflect.Ptr:
			innerType := f.Type().Elem()
			if innerType.Kind() == reflect.Slice {
				continue
			}
			if innerType.Kind() == reflect.Struct {
				new := checkStruct(tag, f.Elem())
				if new == nil {
					continue
				} else if len(new) > 0 {
					keys = append(keys, new...)
					continue
				}
			}
		case reflect.Struct:
			new := checkStruct(tag, f)
			if new == nil {
				continue
			} else if len(new) > 0 {
				keys = append(keys, new...)
				continue
			}
		}
		keys = append(keys, key)
	}
	return keys
}

func checkStruct(tag string, val reflect.Value) []string {
	if !val.Type().Implements(reflect.TypeOf((*driver.Valuer)(nil)).Elem()) {
		if val.Type() != reflect.TypeOf(time.Time{}) {
			if tag != "" {
				return nil
			}
			return getStructureKeys(val.Interface())
		}
	}
	return []string{}
}

func createTable(table, columns string) string {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(", table)
	if config.Env.DBDriver == "mysql" {
		sql += "id INT PRIMARY KEY AUTO_INCREMENT, "
	} else if config.Env.DBDriver == "sqlite3" {
		sql += "id INTEGER PRIMARY KEY AUTOINCREMENT, "
	} else if config.Env.DBDriver == "postgres" {
		sql += "id SERIAL PRIMARY KEY, "
	} else {
		sql += "id INTEGER PRIMARY KEY, " // Unknown SQL databse, defulting to just making id the primary key
	}
	sql += "created_at DATETIME NOT NULL DEFAULT " + domain.DBNow() + ", " // TODO: Test for SQLite
	sql += "updated_at DATETIME NOT NULL DEFAULT " + domain.DBNow() + ", "
	sql += "deleted_at DATETIME, "
	sql += columns
	sql += ")"
	return sql
}
