package domain

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/minebarteksa/clean-show/config"
)

type H map[string]interface{}

type DBModel struct {
	ID        uint         `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type Stmt interface {
	Get(dst interface{}, args interface{}) error
	Exec(args interface{}) (sql.Result, error)
	Select(dst interface{}, args interface{}) error
}

type Tx interface {
	Stmt(Stmt) Stmt
}

type DB interface {
	Prepare(query string) Stmt
	PrepareInsertStruct(table string, arg any) Stmt
	PrepareSelect(table, where string) Stmt
	PrepareUpdate(table, set, where string) Stmt
	PrepareUpdateStruct(table string, arg any, where string) Stmt
	PrepareSoftDelete(table, where string) Stmt

	Exec(query string, args ...any) (sql.Result, error)
	Select(dst interface{}, query string, args ...interface{}) error

	Transaction(func(tx Tx) error) error
}

type FromString interface {
	FromString(string)
}

type DBArray[T any] []T

func (a *DBArray[T]) Scan(val any) error {
	var split []string
	switch v := val.(type) {
	case []byte:
		split = strings.Split(string(v), ";")
	case string:
		split = strings.Split(v, ";")
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	for _, val := range split {
		*a = append(*a, fromString[T](val))
	}
	return nil
}

func fromString[T any](in string) T {
	var a T
	val := reflect.ValueOf(&a)
	if val.Type().Implements(reflect.TypeOf((*FromString)(nil)).Elem()) {
		val.MethodByName("FromString").Call([]reflect.Value{reflect.ValueOf(in)})
		return a
	} else if val.Elem().Kind() == reflect.String {
		val.Elem().SetString(in)
		return a
	}
	log.Panicf("type %T does not implement FromString", a)
	return a
}

func (a DBArray[any]) Value() (driver.Value, error) {
	var tmp []string
	for _, val := range a {
		tmp = append(tmp, fmt.Sprint(val))
	}
	return strings.Join(tmp, ";"), nil
}

func DBNow() string {
	if config.Env.DBDriver == "sqlite3" {
		return "datetime('now', 'localhost')"
	}
	return "NOW()"
}

func DBInterval(base string, durations ...time.Duration) string {
	var durs []string
	for _, duration := range durations {
		hours := int(duration.Hours())
		if hours != 0 {
			durs = append(durs, fmt.Sprintf("%d HOUR", hours))
		}
		minutes := int(duration.Minutes()) % 60
		if minutes != 0 {
			durs = append(durs, fmt.Sprintf("%d MINUTE", minutes))
		}
		seconds := int(duration.Seconds()) % 60
		if seconds != 0 {
			durs = append(durs, fmt.Sprintf("%d SECOND", seconds))
		}
	}
	if config.Env.DBDriver == "sqlite3" {
		return lite(base, durs)
	} else if config.Env.DBDriver == "postgres" {
		return post(base, durs)
	}
	return my(base, durs)
}

func lite(base string, durations []string) string {
	interval := "datetime(" + base + ", "
	for _, dur := range durations {
		interval += "'" + dur + "', "
	}
	return interval[:len(interval)-2] + ")"
}

func post(base string, durations []string) string {
	interval := base
	for _, dur := range durations {
		interval += " + INTERVAL '" + dur + "'"
	}
	return interval
}

func my(base string, durations []string) string {
	interval := base
	for _, dur := range durations {
		interval += " + INTERVAL " + dur
	}
	return interval
}
