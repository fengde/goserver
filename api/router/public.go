package router

import (
	"goserver/api/handler"

	"github.com/gin-gonic/gin"
)

func public(r *gin.Engine) {

	g := r.Group("/")

	// 用于自动化运维健康检查
	GET(g, "/health", handler.Health)
	// demo用于日常验证
	POST(g, "/demo", handler.Demo)
	// 用户登录
	POST(g, "/api/user/login", handler.Login)
}
