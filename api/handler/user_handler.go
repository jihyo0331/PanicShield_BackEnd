package handler

import (
	"net/http"
	"ps_backend/db"
	"ps_backend/model"
	"strconv"

	jwt "ps_backend/pkg/middleware"

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
	// Generate JWT tokens
	accessToken, refreshToken, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "토큰 생성 실패"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "회원가입 성공",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
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
	// Generate JWT tokens
	accessToken, refreshToken, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "토큰 생성 실패"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":        "로그인 성공",
		"access_token":   accessToken,
		"refresh_token":  refreshToken,
		"user_id":        user.ID,
		"speaking_style": user.SpeakingStyle,
		"tone":           user.Tone,
	})
}

// GetProfile retrieves a user's information by query parameter user_id.
func GetProfile(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "유효하지 않은 user_id"})
		return
	}

	var user model.User
	if err := database.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "사용자 정보 조회 실패"})
		return
	}

	// 반환 시 비밀번호 해시 제외
	response := gin.H{
		"id":             user.ID,
		"username":       user.Username,
		"phone_number":   user.PhoneNumber,
		"speaking_style": user.SpeakingStyle,
		"tone":           user.Tone,
		"verified":       user.Verified,
		"created_at":     user.CreatedAt,
	}
	c.JSON(http.StatusOK, response)
}
