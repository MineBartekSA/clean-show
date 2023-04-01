package router

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/minebarteksa/clean-show/domain"
)

type customContext struct {
	*gin.Context
	decoder *json.Decoder
}

func NewContext(c *gin.Context) domain.Context {
	var d *json.Decoder
	if c.Request.Method != "GET" {
		d = json.NewDecoder(c.Request.Body)
	}
	return &customContext{c, d}
}

func (cc *customContext) UnmarshalBody(in any) error {
	if cc.decoder != nil {
		return cc.decoder.Decode(in)
	}
	return fmt.Errorf("no decoder")
}
