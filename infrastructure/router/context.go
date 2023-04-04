package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
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

func (cc *customContext) Error(err error) {
	if e, ok := err.(*domain.Error); ok {
		status := e.Send(cc)
		data := []interface{}{}
		if e.Caller != "" {
			data = append(data, []interface{}{"caller", e.Caller}...)
		}
		if e.Source != nil {
			data = append(data, []interface{}{"source", e.Source}...)
		}
		if status < 500 {
			Log.Debugw(e.Message, data...)
		} else {
			Log.Errorw(e.Message, data...)
		}
	} else {
		Log.Errorw("unexpected internal error", "url", cc.Request.URL, "err", err)
		cc.JSON(http.StatusInternalServerError, domain.H{"error": "internal server error"})
	}
}
