package router

import (
	"goserver/api/handler"
	"goserver/api/middleware"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	r.Use(middleware.Statistics())
	r.Use(middleware.AccessLog())

	// 用于自动化运维健康检查
	GET(&r.RouterGroup, "/health", handler.Health)
	// demo用于日常验证
	POST(&r.RouterGroup, "/demo", handler.Demo)

	public := r.Group("/")
	POST(public, "/api/user/login", handler.Login)

	private := r.Group("/")
	private.Use(middleware.Jwt())
	POST(private, "/api/user/info", handler.Info)
}
