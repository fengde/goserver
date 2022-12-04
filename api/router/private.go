package router

import (
	"goserver/api/handler"
	"goserver/api/middleware"

	"github.com/gin-gonic/gin"
)

func private(r *gin.Engine) {
	g := r.Group("/")
	g.Use(middleware.Jwt())

	POST(g, "/api/user/info", handler.Info)
}
