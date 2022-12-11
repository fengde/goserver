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
	// 用户信息查询
	POST(g, "/api/user/info", handler.Info)
	// 所有权限组查询
	POST(g, "/api/permission/groups", handler.GetPermissionGroups)
}
