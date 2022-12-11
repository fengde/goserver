package router

import (
	"goserver/api/handler"
	"goserver/api/middleware"

	"github.com/gin-gonic/gin"
)

// 登录，并且需要授权的url在这里注册
func privateRbac(r *gin.Engine) {
	g := r.Group("/")

	g.Use(middleware.Jwt()).Use(middleware.Rbac())

	POST(g, "/api/role/new", handler.NewRole)
	POST(g, "/api/role/delete", handler.DeleteRole)
}
