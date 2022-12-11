package middleware

import (
	"goserver/api/handler"
	"goserver/service/serviceRbac"
	"goserver/service/serviceUser"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Rbac() gin.HandlerFunc {
	return func(ginc *gin.Context) {
		ctx := handler.GetCtx(ginc)
		userId := handler.GetUserId(ginc)

		if !serviceUser.IsSuper(ctx, userId) && !serviceRbac.Check(ctx, userId, ginc.Request.URL.RequestURI(), ginc.Request.Method) {
			ginc.AbortWithStatusJSON(http.StatusOK, gin.H{
				"status":  "failed",
				"message": "access without permission",
				"data":    map[string]interface{}{},
			})
			return
		}

		ginc.Next()
	}
}
