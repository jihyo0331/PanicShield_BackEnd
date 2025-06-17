package api

import (
	"ps_backend/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/signup", handler.SignUp)
		api.POST("/login", handler.Login)
		api.POST("/chat", handler.ChatWithGemini)
	}

	return r
}
