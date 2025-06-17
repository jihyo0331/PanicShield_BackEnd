package handler

import (
	"net/http"
	"ps_backend/db"
	"ps_backend/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var database = db.GetDB()

type SignUpRequest struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	PhoneNumber   string `json:"phone_number"`
	SpeakingStyle string `json:"speaking_style"` // "반말" 등
	Tone          string `json:"tone"`           // "유머" 등
}

func SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "데이터 형식 오류"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "비밀번호 해시 실패"})
		return
	}

	user := model.User{
		Username:      req.Username,
		PasswordHash:  string(hash),
		PhoneNumber:   req.PhoneNumber,
		SpeakingStyle: req.SpeakingStyle,
		Tone:          req.Tone,
	}
	if err := database.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "유저 생성 실패(중복?)"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "회원가입 성공"})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "데이터 형식 오류"})
		return
	}
	var user model.User
	if err := database.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "존재하지 않는 유저"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "비밀번호 불일치"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "로그인 성공", "user_id": user.ID, "speaking_style": user.SpeakingStyle, "tone": user.Tone})
}
