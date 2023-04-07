package test

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/minebarteksa/clean-show/domain"
)

type MockContext struct {
	ParamMap  map[string]string
	QueryMap  map[string]string
	HeaderMap map[string]string
	In        any

	OutCookie map[string]string
	OutStatus int
	OutError  error
	Out       any
}

func (mc *MockContext) Param(key string) string {
	if val, ok := mc.ParamMap[key]; ok {
		return val
	}
	return ""
}

func (mc *MockContext) Query(key string) string {
	if val, ok := mc.QueryMap[key]; ok {
		return val
	}
	return ""
}

func (mc *MockContext) GetHeader(header string) string {
	if val, ok := mc.HeaderMap[header]; ok {
		return val
	}
	return ""
}

func (mc *MockContext) UnmarshalBody(in any) error {
	val := reflect.ValueOf(in)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}
	ival := reflect.ValueOf(mc.In)
	if ival.Kind() == reflect.Ptr {
		ival = ival.Elem()
	}
	val.Elem().Set(ival)
	return nil
}

func (mc *MockContext) Cookie(name string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (mc *MockContext) Error(err error) {
	mc.OutStatus = http.StatusInternalServerError
	mc.OutError = err
	mc.Out = domain.H{"error": "internal server error"}
	if e, ok := err.(*domain.Error); ok {
		e.Send(mc)
	}
}

func (mc *MockContext) JSON(code int, i any) {
	mc.OutStatus = code
	mc.Out = i
}

func (mc *MockContext) Status(code int) {
	mc.OutStatus = code
	mc.Out = nil
}

func (mc *MockContext) String(code int, format string, value ...any) {
	mc.OutStatus = code
	mc.Out = fmt.Sprintf(format, value...)
}

func (mc *MockContext) SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) {
	mc.OutCookie[name] = value
}
