package domain

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"runtime/debug"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
)

var callerSkip int

func init() {
	_, b, _, _ := runtime.Caller(0)
	callerSkip = len(filepath.Dir(filepath.Dir(b)))
}

type Error struct {
	Source error
	Stack  []byte
	Caller string

	Message string
	Status  int
}

func (e *Error) Error() string {
	if e.Source == nil {
		return e.Message
	} else {
		return fmt.Sprintf("%s: %s", e.Message, e.Source)
	}
}

func (e *Error) Call() *Error {
	return e.SetCall(2)
}

func (e *Error) SetCall(skip int) *Error {
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		e.Caller = fmt.Sprintf("%s:%d", file[callerSkip:], line)
	}
	return e
}

func (e *Error) Wrap(wrap error) *Error {
	e.Source = wrap
	return e
}

func (e *Error) SetStack() *Error {
	e.Stack = debug.Stack()
	return e
}

func (e *Error) Send(context Context) int {
	if e.Status != 0 {
		context.JSON(e.Status, H{"error": e.Message})
		return e.Status
	}
	if err, ok := e.Source.(*Error); ok {
		return err.Send(context)
	}
	context.JSON(http.StatusInternalServerError, H{"error": "internal server error"})
	return http.StatusInternalServerError
}

func Fatal(source error, message string) *Error {
	return &Error{
		Source:  source,
		Stack:   debug.Stack(),
		Message: message,
	}
}

func Check(err error, expect *Error) bool {
	if e, ok := err.(*Error); ok {
		return e.Message == expect.Message
	}
	return false
}

var (
	ErrUnauthorized = &Error{
		Message: "unauthorized",
		Status:  http.StatusUnauthorized,
	}
	ErrBadRequest = &Error{
		Message: "bad request",
		Status:  http.StatusBadRequest,
	}
	ErrBadPassword = &Error{
		Message: "invalid password",
		Status:  http.StatusBadRequest,
	}
	ErrNotFound = &Error{
		Message: "resource not found",
		Status:  http.StatusNotFound,
	}
	ErrDuplicate = &Error{
		Message: "duplicate resource",
		Status:  http.StatusBadRequest,
	}
	ErrNullValue = &Error{
		Message: "null value inserted into a not-null column",
		Status:  http.StatusBadRequest,
	}
)

func SQLError(source error) error {
	if source == nil {
		return nil
	}
	if source == sql.ErrNoRows {
		return ErrNotFound.Wrap(source).SetCall(2)
	} else if source == sql.ErrTxDone {
		return Fatal(source, "transaction terminated").SetCall(2)
	}
	switch err := source.(type) {
	case *pq.Error:
		if err.Code == "23505" {
			return ErrDuplicate.Wrap(source).SetCall(2)
		} else if err.Code == "23502" {
			return ErrNullValue.Wrap(source).SetCall(2)
		}
	case *mysql.MySQLError:
		if err.Number == 1062 || err.Number == 1586 {
			return ErrDuplicate.Wrap(source).SetCall(2)
		} else if err.Number == 1048 {
			return ErrNullValue.Wrap(source).SetCall(2)
		}
	case *sqlite3.Error:
		if err.ExtendedCode == sqlite3.ErrConstraintUnique {
			return ErrDuplicate.Wrap(source).SetCall(2)
		} else if err.ExtendedCode == sqlite3.ErrConstraintNotNull {
			return ErrNullValue.Wrap(source).SetCall(2)
		}
	}
	return Fatal(source, "unknown sql error")
}
