package handler

import (
	"net/http"
	"ps_backend/db"
	dto "ps_backend/dto"
	panicService "ps_backend/internal/panic_guide"

	"github.com/gin-gonic/gin"
)

func ListPanicGuides(c *gin.Context) {
	guides, err := panicService.NewService(db.GetDB()).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve panic guides"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": guides})
}

func AddPanicGuide(c *gin.Context) {
	var req dto.PanicGuideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	_, err := panicService.NewService(db.GetDB()).Create(req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add panic guide"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Panic guide added successfully"})
}

func BookmarkPanicGuide(c *gin.Context) {
	var req dto.BookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	err := panicService.NewService(db.GetDB()).Bookmark(req.UserID, req.PanicGuideID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bookmark panic guide"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Panic guide bookmarked successfully"})
}

func ListUserBookmarks(c *gin.Context) {
	var req dto.BookmarkRequest

	// Try to bind JSON first
	if err := c.ShouldBindJSON(&req); err != nil || req.UserID == 0 {
		// If JSON bind fails or userID is zero, try to get user_id from query param
		userIDQuery := c.Query("user_id")
		if userIDQuery == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "user_id is required"})
			return
		}
		var err error
		var parsed uint64
		parsed, err = dto.ParseUint(userIDQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user_id"})
			return
		}
		req.UserID = uint(parsed)
		if req.UserID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user_id"})
			return
		}
	}

	bookmarkedGuides, err := panicService.NewService(db.GetDB()).GetBookmarks(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve bookmarked panic guides"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookmarkedGuides})
}
