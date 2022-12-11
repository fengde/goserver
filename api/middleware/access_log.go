package middleware

import (
	"bytes"
	"goserver/api/handler"
	"io/ioutil"
	"strings"
	"time"

	"github.com/fengde/gocommon/logx"
	"github.com/gin-gonic/gin"
)

func AccessLog() gin.HandlerFunc {
	return func(ginc *gin.Context) {
		start := time.Now()

		ctx := handler.GetCtx(ginc)

		// start
		{
			body, _ := ginc.GetRawData()
			// 将原body塞回去
			ginc.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			headers := []string{}

			for k, v := range ginc.Request.Header {
				headers = append(headers, k+":"+v[0])
			}

			// 请求前
			logx.InfofWithCtx(ctx, `收到请求 | %v | %v%v | %v | From: %v | Header: %v | Body: %v | Body size: %v bytes`,
				ginc.Request.Method,
				ginc.Request.Host,
				ginc.Request.URL,
				ginc.Request.Proto,
				ginc.RemoteIP(),
				strings.Join(headers, "\n"),
				string(body),
				len(body))
		}

		ginc.Next()

		// end
		{
			headers := []string{}

			for k, v := range ginc.Writer.Header() {
				headers = append(headers, k+":"+v[0])
			}

			out := ginc.GetString("out")

			logx.InfofWithCtx(ctx, `请求结束 | http status: %v | Header: %v | Body: %v | Body size: %v bytes | 请求耗时: %v`,
				ginc.Writer.Status(), strings.Join(headers, "\n"), out, len(out), time.Since(start))

		}
	}
}
