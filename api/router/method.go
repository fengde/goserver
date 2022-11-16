package router

import (
	"goserver/api/handler"

	"github.com/gin-gonic/gin"
)

func GET(group *gin.RouterGroup, relativePath string, f interface{}) {
	group.GET(relativePath, handler.WrapF(f))
}

func POST(group *gin.RouterGroup, relativePath string, f interface{}) {
	group.POST(relativePath, handler.WrapF(f))
}

func PUT(group *gin.RouterGroup, relativePath string, f interface{}) {
	group.PUT(relativePath, handler.WrapF(f))
}

func DELETE(group *gin.RouterGroup, relativePath string, f interface{}) {
	group.DELETE(relativePath, handler.WrapF(f))
}

func HEAD(group *gin.RouterGroup, relativePath string, f interface{}) {
	group.HEAD(relativePath, handler.WrapF(f))
}

func Any(group *gin.RouterGroup, relativePath string, f interface{}) {
	group.Any(relativePath, handler.WrapF(f))
}
