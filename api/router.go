package api

import (
	"ps_backend/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.SignUp)
			auth.POST("/signin", handler.Login)
			auth.POST("/refresh", handler.RefreshToken)
		}

		api.GET("/users/me", handler.GetProfile)

		// Interests
		api.POST("/interests", handler.AddInterest)
		api.GET("/interests", handler.ListInterests)
		api.POST("/vitals", handler.RegisterVital)
		api.GET("/vitals", handler.ListVitals)

		// Chatbot
		api.GET("/chat", handler.GetChatHistory)
		api.POST("/chat", handler.ChatWithGemini)

		// Panic guides
		api.GET("/panic-guides", handler.ListPanicGuides)
		api.POST("/panic-guides", handler.AddPanicGuide)
		api.POST("/panic-guides/bookmark", handler.BookmarkPanicGuide)
		api.GET("/panic-guides/bookmarks", handler.ListUserBookmarks)
	}

	return r
}
