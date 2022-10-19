package middleware

import (
	"server/service/serviceHttp"

	"github.com/gin-gonic/gin"
)

// 统计分析
func Statistics() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceHttp.RequestIn(c)
		defer serviceHttp.RequestOut(c)
		c.Next()
	}
}
