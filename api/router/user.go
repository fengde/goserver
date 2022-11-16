package router

import (
	"goserver/api/handler"
	"goserver/api/middleware"

	"github.com/gin-gonic/gin"
)

func user(r *gin.Engine) {
	public := r.Group("/")
	public.POST("/api/user/login", handler.WrapF(handler.Login))

	private := r.Group("/")
	private.Use(middleware.Jwt())
	private.POST("/api/user/info", handler.WrapF(handler.Info))
}
