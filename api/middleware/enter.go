package middleware

import (
	"fmt"

	"github.com/fengde/gocommon/logx"
	"github.com/fengde/gocommon/timex"
	"github.com/fengde/gocommon/toolx"
	"github.com/gin-gonic/gin"
)

var XRequestIdHeader = "x-request-id"

func Enter() gin.HandlerFunc {
	return func(c *gin.Context) {
		xRequestId := c.GetHeader(XRequestIdHeader)
		if xRequestId == "" {
			xRequestId = fmt.Sprintf("%v%s", timex.NowUnixNano(), toolx.NewNumberCode(4))
		}

		c.Set(XRequestIdHeader, xRequestId)
		c.Set("ctx", logx.NewCtx(xRequestId))
		c.Next()
	}
}
