package users

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/users", handler.Create)
	router.GET("/users", handler.List)
}
