package handler

import (
	"net/http"
	"time"

	"ps_backend/db"
	"ps_backend/dto"
	vitalService "ps_backend/internal/vital"
	"ps_backend/model"

	"github.com/gin-gonic/gin"
)

var vitalSvc = vitalService.NewService(db.GetDB())

func RegisterVital(c *gin.Context) {
	var req dto.VitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	entry := &model.VitalSign{
		UserID:      req.UserID,
		HeartRate:   req.HeartRate,
		BreathRate:  req.BreathRate,
		StressLevel: req.StressLevel,
		MeasuredAt:  time.Now(),
	}
	if err := vitalSvc.CreateVital(entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vital record"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Vital record created successfully"})
}

func ListVitals(c *gin.Context) {
	var userID uint
	var err error

	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		var req struct {
			UserID string `json:"user_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}
		userIDStr = req.UserID
	}

	parsed, err := dto.ParseUint(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	userID = uint(parsed)

	vitals, err := vitalSvc.GetVitalsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vital records"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": vitals})
}
