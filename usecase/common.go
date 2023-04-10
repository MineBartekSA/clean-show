package usecase

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/minebarteksa/clean-show/domain"
)

func PatchModel(model any, data map[string]any) (err error) {
	val := reflect.ValueOf(model)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct")
	}
	dbm := reflect.TypeOf(domain.DBModel{})
	if val.Type() == dbm {
		return fmt.Errorf("DBModel can not be patched")
	}
	defer func() {
		if r := recover(); r != nil {
			err = domain.ErrBadRequest.Wrap(err).Call()
			return
		}
	}()
	max := val.NumField()
	for i := 0; i < max; i += 1 {
		t := val.Type().Field(i)
		if t.Type == dbm {
			continue
		}
		key := t.Name
		if tag, ok := t.Tag.Lookup("json"); ok {
			new := strings.SplitN(tag, ",", 2)[0]
			if new != "" {
				key = new
			}
		}
		if tag, ok := t.Tag.Lookup("patch"); ok {
			if tag == "-" {
				continue
			}
		}
		if key == "-" {
			continue
		}
		if value, ok := data[key]; ok {
			f := val.Field(i)
			v := reflect.ValueOf(value)
			switch f.Kind() {
			case reflect.Int:
				switch v.Kind() {
				case reflect.Int64, reflect.Int:
					f.Set(v)
				case reflect.Float64:
					f.SetInt(int64(value.(float64)))
				default:
					return fmt.Errorf("can not put type %T into %s", value, f.Type())
				}
			case reflect.Float64:
				switch v.Kind() {
				case reflect.Float64:
					f.Set(v)
				case reflect.Int64:
					f.SetFloat(float64(value.(int64)))
				default:
					return fmt.Errorf("can not put type %T into %s", value, f.Type())
				}
			default:
				if v.Kind() == reflect.Slice {
					new := reflect.New(f.Type())
					if anySlice := new.MethodByName("FromInterfaceSlice"); anySlice != reflect.ValueOf(nil) {
						anySlice.Call([]reflect.Value{v})
					}
					f.Set(new.Elem())
				} else {
					f.Set(v)
				}
			}
		}
	}
	return
}

func GetLimitPage(context domain.Context) (int, int) {
	limit := 10
	page := 1
	if query := context.Query("limit"); query != "" {
		l, err := strconv.ParseInt(query, 10, 64)
		if err == nil {
			limit = int(l)
		}
	}
	if query := context.Query("page"); query != "" {
		p, err := strconv.ParseInt(query, 10, 64)
		if err == nil {
			page = int(p)
		}
	}
	return limit, page
}
