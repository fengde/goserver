package router

import (
	"goserver/api/middleware"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	// 中间件设置
	r.Use(middleware.Enter()).Use(middleware.AccessLog()).Use(middleware.Statistics())
	// 公开路由设置
	public(r)
	// 鉴权路由设置
	private(r)
}
