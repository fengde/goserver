package router

import (
	"goserver/api/handler"
	"goserver/api/middleware"

	"github.com/gin-gonic/gin"
)

// 登录即可访问在这里注册
func private(r *gin.Engine) {
	g := r.Group("/")

	g.Use(middleware.Jwt())

	POST(g, "/api/user/info", handler.Info)
}
