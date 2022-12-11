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
	// 登录即可访问
	private(r)
	// 登录还需授权方可访问
	privateRbac(r)
}
