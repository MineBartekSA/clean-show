package database

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
)

type sqldb struct {
	*sqlx.DB
}

func NewSqlDB(db *sqlx.DB) domain.DB {
	createTable(db, "audit_log", `
		type INT NOT NULL,
		resource_type INT NOT NULL,
		resource_id INT NOT NULL,
		executor_id INT NOT NULL
	`)
	createTable(db, "sessions", `
		account_id INT NOT NULL,
		token VARCHAR(128) NOT NULL
	`, "updated_at token", "UNIQUE token")
	createTable(db, "accounts", `
		type INT NOT NULL,
		email VARCHAR(320) NOT NULL,
		hash TEXT NOT NULL,
		name TEXT NOT NULL,
		surname TEXT NOT NULL
	`, "UNIQUE email") // TODO: Change hash VARCHAR // namd and surname should also be VARCHAR
	createTable(db, "products", `
		status INT NOT NULL,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		price REAL NOT NULL,
		images TEXT NOT NULL
	`) // TODO: name should be chnaged to VARCHAR
	createTable(db, "orders", `
		status INT NOT NULL,
		order_by INT NOT NULL,
		shipping_address TEXT NOT NULL,
		invoice_address TEXT NOT NULL,
		products TEXT NOT NULL,
		shipping_price REAL NOT NULL,
		total REAL NOT NULL
	`, "order_by")

	return &sqldb{db}
}

func (db *sqldb) PrepareInsertStruct(table string, arg any) domain.Stmt {
	cols := getStructureKeys(arg)
	sql := ""
	if config.Env.DBDriver == "mysql" {
		sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (:%s)", table, strings.Join(cols, ", "), strings.Join(cols, ", :"))
	} else {
		sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (:%s) RETURNING id", table, strings.Join(cols, ", "), strings.Join(cols, ", :"))
	}
	stmt, err := db.innerPrepare(sql)
	if err != nil {
		Log.Panicw("failed to prepare a named insert statement", "statement", sql, "err", err)
	}
	return stmt
}

func (db *sqldb) PrepareSelect(table string, where string) domain.Stmt {
	if where != "" {
		where += " AND "
	}
	where += "deleted_at IS NULL"
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, where)
	stmt, err := db.innerPrepare(sql)
	if err != nil {
		Log.Panicw("failed to prepare a named select statement", "statement", sql, "err", err)
	}
	return stmt
}

func (db *sqldb) PrepareUpdateStruct(table string, arg any, where string) domain.Stmt {
	cols := getStructureKeys(arg)
	set := ""
	for _, col := range cols {
		set += fmt.Sprintf("%s = :%s, ", col, col)
	}
	return db.PrepareUpdate(table, set[:len(set)-2], where)
}

func (db *sqldb) PrepareUpdate(table, set, where string) domain.Stmt {
	if set != "" {
		set += ", "
	}
	set += "updated_at = " + domain.DBNow()
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, set, where)
	stmt, err := db.innerPrepare(sql)
	if err != nil {
		Log.Panicw("failed to prepare a named update statement", "statement", sql, "err", err)
	}
	return stmt
}

func (db *sqldb) PrepareSoftDelete(table string, where string) domain.Stmt {
	sql := fmt.Sprintf("UPDATE %s SET deleted_at = %s WHERE %s", table, domain.DBNow(), where)
	stmt, err := db.innerPrepare(sql)
	if err != nil {
		Log.Panicw("failed to prepare a named soft delete statement", "statement", sql, "err", err)
	}
	return stmt
}

func (db *sqldb) Prepare(query string) domain.Stmt {
	stmt, err := db.innerPrepare(query)
	if err != nil {
		Log.Panicw("failed to prepare a named statement", "statement", query, "err", err)
	}
	return stmt
}

func (db *sqldb) innerPrepare(query string) (domain.Stmt, error) {
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

func getStructureKeys(arg any) []string {
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

type index struct {
	table   string
	unique  bool
	columns []string
}

func (i *index) Name() string {
	name := ""
	if i.unique {
		name += "u"
	}
	return name + fmt.Sprintf("idx_%s_%s", i.table, strings.Join(i.columns, "_"))
}

func (i *index) String() string {
	sql := "CREATE "
	if i.unique {
		sql += "UNIQUE "
	}
	sql += "INDEX "
	if config.Env.DBDriver != "mysql" {
		sql += "IF NOT EXISTS "
	}
	sql += i.Name() + " "
	sql += "ON " + i.table + " "
	sql += "(" + strings.Join(i.columns, ", ") + ")"
	return sql
}

func indexFromString(table, in string) *index {
	split := strings.Split(in, " ")
	i := index{table: table}
	if strings.ToLower(split[0]) == "unique" {
		i.unique = true
		if len(split) == 1 {
			Log.Panicw("invalif index format", "table", table, "raw", in)
		}
		split = split[1:]
	}
	i.columns = split
	return &i
}

func createTable(db *sqlx.DB, table, columns string, indexes ...string) {
	now := domain.DBNow()
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(", table)
	if config.Env.DBDriver == "mysql" {
		sql += "id INT PRIMARY KEY AUTO_INCREMENT, "
	} else if config.Env.DBDriver == "sqlite3" {
		sql += "id INTEGER PRIMARY KEY AUTOINCREMENT, "
		now = "(" + now + ")"
	} else if config.Env.DBDriver == "postgres" {
		sql += "id SERIAL PRIMARY KEY, "
	} else {
		sql += "id INTEGER PRIMARY KEY, " // Unknown SQL databse, defulting to just making id the primary key
	}
	sql += "created_at DATETIME NOT NULL DEFAULT " + now + ", "
	sql += "updated_at DATETIME NOT NULL DEFAULT " + now + ", "
	sql += "deleted_at DATETIME, "
	sql += columns
	sql += ")"
	db.MustExec(sql)
	bindex := []*index{indexFromString(table, "deleted_at")}
	for _, raw := range indexes {
		bindex = append(bindex, indexFromString(table, raw))
	}
	for _, i := range bindex {
		if config.Env.DBDriver == "mysql" {
			var res uint
			err := db.Get(&res, fmt.Sprintf("SELECT COUNT(1) FROM INFORMATION_SCHEMA.STATISTICS WHERE table_schema=DATABASE() AND table_name='%s' AND index_name='%s'", table, i.Name()))
			if err != nil {
				Log.Panicw("failed to check for index", "table", table, "index", i.Name(), "err", err)
			}
			if res != 0 {
				continue
			}
		}
		Log.Debugln(i.String())
		db.MustExec(i.String())
	}
}
