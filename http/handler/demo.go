package handler

import (
	"github.com/fengde/gocommon/jsonx"
	"github.com/fengde/gocommon/logx"
)

type DemoRequest struct {
	User string `form:"user"` // uri参数
	Name string `json:"name"` // json参数
	Age  string `form:"age"`  // form参数
}

func Demo(c *Context, r *DemoRequest) {
	logx.Info(jsonx.MarshalToStringNoErr(r))
	c.OutSuccess(nil)
}
