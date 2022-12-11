package router

import (
	"goserver/api/handler"
	"goserver/api/middleware"

	"github.com/gin-gonic/gin"
)

func private(r *gin.Engine) {
	g := r.Group("/")

	g.Use(middleware.Jwt()).Use(middleware.Rbac())

	POST(g, "/api/user/info", handler.Info)
	POST(g, "/api/role/new", handler.NewRole)
	POST(g, "/api/role/delete", handler.DeleteRole)
}
