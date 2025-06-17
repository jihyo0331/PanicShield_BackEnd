package handler

import (
	"net/http"
	"time"

	"ps_backend/db"
	dto "ps_backend/dto"
	authservice "ps_backend/internal/auth"
	"ps_backend/model"
	jwt "ps_backend/pkg/middleware"
	"ps_backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

var authSvc = authservice.NewAuthService(db.GetDB())

// Register handles user registration requests
func Register(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Check for existing username or phone
	if exists, err := authSvc.UserExistsByUsername(req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check username"})
		return
	} else if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	if exists, err := authSvc.UserExistsByPhone(req.PhoneNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check phone number"})
		return
	} else if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number already registered"})
		return
	}

	// Hash password
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := &model.User{
		Username:      req.Username,
		PasswordHash:  hashed,
		PhoneNumber:   req.PhoneNumber,
		SpeakingStyle: req.SpeakingStyle,
		Tone:          req.Tone,
		CreatedAt:     time.Now(),
	}

	if err := authSvc.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered", "user_id": user.ID})
}

// SignIn handles authentication and JWT issuance
func SignIn(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	user, err := authSvc.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, refreshToken, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// VerifyPhone handles phone verification using OTP codes
func VerifyPhone(c *gin.Context) {
	var req dto.VerifyPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate OTP code
	if !authSvc.ValidateOTP(req.UserID, req.Code) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired verification code"})
		return
	}

	// Mark user as verified
	if err := authSvc.MarkPhoneVerified(req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify phone"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Phone number verified"})
}

// RefreshToken issues a new JWT using a valid refresh token
func RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	accessToken, refreshToken, err := jwt.GenerateTokenFromRefresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Token refreshed",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
