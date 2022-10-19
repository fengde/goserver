package router

import (
	"server/http/handler"
	"server/http/middleware"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	r.Use(middleware.Statistics())
	r.Use(middleware.AccessLog())

	// 用于自动化运维健康检查
	r.GET("/health", handler.WrapF(handler.Health))
	// demo用于日常验证
	r.POST("/demo", handler.WrapF(handler.Demo))

	user(r)
}
