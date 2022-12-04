package router

import (
	"goserver/api/middleware"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	// 中间件设置
	r.Use(middleware.Enter())
	r.Use(middleware.Jaeger())
	r.Use(middleware.Statistics())
	r.Use(middleware.AccessLog())
	// 路由设置
	public(r)
	private(r)
}
