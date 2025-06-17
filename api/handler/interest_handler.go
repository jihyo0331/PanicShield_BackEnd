package handler

import (
	"net/http"
	"ps_backend/db"
	"ps_backend/model"

	"github.com/gin-gonic/gin"
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

	// 관심사 없으면 생성, 있으면 그대로 사용
	var interest model.Interest
	if err := databaseI.Where("name = ?", req.Interest).First(&interest).Error; err != nil {
		interest = model.Interest{Name: req.Interest}
		databaseI.Create(&interest)
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
