package middleware

import (
	"goserver/service/serviceJwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Jwt检查校验
func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"status":  "login",
				"message": "token is missing",
				"data":    map[string]interface{}{},
			})
			return
		}

		jwtClaims, err := serviceJwt.ParseJwt(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"status":  "login",
				"message": err.Error(),
				"data":    map[string]interface{}{},
			})
			return
		}

		c.Set("user_id", jwtClaims.UserId)
		// 处理请求
		c.Next()
	}
}
