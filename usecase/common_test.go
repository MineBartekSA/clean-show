package usecase_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/infrastructure/router"
	"github.com/minebarteksa/clean-show/usecase"
	"github.com/stretchr/testify/assert"
)

type PatchModelTest struct {
	domain.DBModel `json:"-"`
	String         string `json:"string"`
	Integer        int
	Float          float64
	Bool           bool
	NoPatch        int `patch:"-"`
}

func TestPatchModel(t *testing.T) {
	var dbm domain.DBModel

	err := usecase.PatchModel(&dbm, map[string]any{})
	assert.Error(t, err)

	err = usecase.PatchModel(dbm, map[string]any{})
	assert.Error(t, err)

	var integer int

	err = usecase.PatchModel(&integer, map[string]any{})
	assert.Error(t, err)

	testStruct := PatchModelTest{
		DBModel: domain.DBModel{
			ID: 100,
		},
		String:  "test",
		Integer: 5,
		Float:   2.7,
		Bool:    false,
		NoPatch: 10,
	}
	err = usecase.PatchModel(&testStruct, map[string]any{
		"ID":      1,
		"string":  "hello",
		"String":  "world",
		"Integer": 300,
		"Float":   float64(5.9),
		"Bool":    true,
		"NoPatch": 9999,
	})
	assert.NoError(t, err)
	assert.Equal(t, uint(100), testStruct.ID)
	assert.Equal(t, "hello", testStruct.String)
	assert.Equal(t, 300, testStruct.Integer)
	assert.Equal(t, float64(5.9), testStruct.Float)
	assert.Equal(t, true, testStruct.Bool)
	assert.Equal(t, 10, testStruct.NoPatch)
}

func TestPatchModelFromJSON(t *testing.T) {
	testStruct := PatchModelTest{}
	var input map[string]any
	err := json.Unmarshal([]byte("{\"ID\":10,\"string\":\"hello there\",\"Integer\":300,\"Float\":5.9,\"Bool\":true}"), &input)
	assert.NoError(t, err)

	err = usecase.PatchModel(&testStruct, input)

	assert.NoError(t, err)
	assert.Equal(t, uint(0), testStruct.ID)
	assert.Equal(t, "hello there", testStruct.String)
	assert.Equal(t, 300, testStruct.Integer)
	assert.Equal(t, float64(5.9), testStruct.Float)
	assert.Equal(t, true, testStruct.Bool)
	assert.Equal(t, 0, testStruct.NoPatch)
}

func TestGetLimitPage(t *testing.T) {
	context := buildContext("localhost")
	limit, page := usecase.GetLimitPage(context)
	assert.Equal(t, 10, limit)
	assert.Equal(t, 1, page)

	context = buildContext("localhost/test?limit=100")
	limit, page = usecase.GetLimitPage(context)
	assert.Equal(t, 100, limit)
	assert.Equal(t, 1, page)

	context = buildContext("localhost/test?page=100")
	limit, page = usecase.GetLimitPage(context)
	assert.Equal(t, 10, limit)
	assert.Equal(t, 100, page)

	context = buildContext("localhost/test?limit=100&page=1000")
	limit, page = usecase.GetLimitPage(context)
	assert.Equal(t, 100, limit)
	assert.Equal(t, 1000, page)
}

func buildContext(rawurl string) domain.Context {
	inner := gin.Context{Request: &http.Request{}}
	inner.Request.URL, _ = url.Parse(rawurl)
	return router.NewContext(&inner)
}
