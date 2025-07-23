package handler

import (
	"errors"
	"net/http"
	"ps_backend/db"
	"ps_backend/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var databaseI = db.GetDB()

type AddInterestRequest struct {
	UserID   uint   `json:"user_id"`
	Interest string `json:"interest"`
}

func AddInterest(c *gin.Context) {
	var req AddInterestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "데이터 형식 오류"})
		return
	}

	var interest model.Interest
	err := databaseI.Where("name = ?", req.Interest).First(&interest).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new interest if not found
		interest = model.Interest{Name: req.Interest}
		if err := databaseI.Create(&interest).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "관심사 생성 실패"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "관심사 조회 실패"})
		return
	}

	userInterest := model.UserInterest{
		UserID:     req.UserID,
		InterestID: interest.ID,
	}
	if err := databaseI.Create(&userInterest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "관심사 추가 실패"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "관심사 등록 성공"})
}

func ListInterests(c *gin.Context) {
	var interests []model.Interest
	if err := databaseI.Find(&interests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "관심사 조회 실패"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"interests": interests})
}
