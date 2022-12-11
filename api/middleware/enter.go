package middleware

import (
	"fmt"

	"github.com/fengde/gocommon/logx"
	"github.com/fengde/gocommon/timex"
	"github.com/fengde/gocommon/toolx"
	"github.com/gin-gonic/gin"
)

var RequestIdHeader = "request-id"

func Enter() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetHeader(RequestIdHeader)
		if requestId == "" {
			requestId = fmt.Sprintf("%v%s", timex.NowUnixNano(), toolx.NewNumberCode(4))
		}

		c.Set("request_id", requestId)
		c.Set("ctx", logx.NewCtx(requestId))
		c.Next()
	}
}
