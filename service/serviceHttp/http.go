package serviceHttp

import (
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type HttpRequestCount struct {
	Uptime            int64 // 系统运行时长
	Total             int64 // 总的请求数
	Total200          int64 // 总的正常请求数
	Current           int64 // 当前正在处理的请求数
	RequestSizeBytes  int64 // 总请求字节数
	ResponseSizeBytes int64 // 总请求返回字节数
}

var count HttpRequestCount

func RequestIn(c *gin.Context) {
	atomic.AddInt64(&count.Total, 1)
	atomic.AddInt64(&count.Current, 1)
	// 统计请求字节数
	size := 0
	if c.Request.URL != nil {
		size = len(c.Request.URL.String())
	}

	size += len(c.Request.Method)
	size += len(c.Request.Proto)
	for name, values := range c.Request.Header {
		size += len(name)
		for _, value := range values {
			size += len(value)
		}
	}
	size += len(c.Request.Host)
	if c.Request.ContentLength != -1 {
		size += int(c.Request.ContentLength)
	}

	atomic.AddInt64(&count.RequestSizeBytes, int64(size))
}

func RequestOut(c *gin.Context) {
	atomic.AddInt64(&count.Current, -1)
	if c.Writer.Status()/100 == 2 {
		atomic.AddInt64(&count.Total200, 1)
	}
	// 统计返回字节数
	size := c.Writer.Size()
	if size < 0 {
		size = 0
	}
	atomic.AddInt64(&count.ResponseSizeBytes, int64(size))
}

// 返回统计信息
func GetCountInfo() HttpRequestCount {
	return count
}
