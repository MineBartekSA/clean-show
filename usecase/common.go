package usecase

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/minebarteksa/clean-show/domain"
)

func PatchModel(model any, data map[string]any) error { // TODO: Test
	val := reflect.ValueOf(model)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct")
	}
	max := val.NumField()
	for i := 0; i < max; i += 1 {
		t := val.Type().Field(i)
		key := strings.ToLower(t.Name)
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
			val.Field(i).Set(reflect.ValueOf(value))
		}
	}
	return nil
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
