package auth

import (
	"errors"
	"os"
	"time"

	"ps_backend/model"

	jwtmw "ps_backend/pkg/middleware"

	jwt "github.com/golang-jwt/jwt/v5"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) UserExistsByUsername(username string) (bool, error) {
	var count int64
	if err := s.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *AuthService) UserExistsByPhone(phone string) (bool, error) {
	var count int64
	if err := s.db.Model(&model.User{}).Where("phone_number = ?", phone).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *AuthService) CreateUser(user *model.User) error {
	return s.db.Create(user).Error
}

func (s *AuthService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) ValidateOTP(userID uint, code string) bool {
	// TODO: Implement OTP validation logic or import the correct package
	return false
}
func (s *AuthService) MarkPhoneVerified(userID uint) error {
	return s.db.Model(&model.User{}).Where("id = ?", userID).Update("phone_verified", true).Error
}

func (s *AuthService) RefreshJWT(refreshToken string) (string, error) {
	refreshSecret := os.Getenv("REFRESH_SECRET")
	if refreshSecret == "" {
		return "", errors.New("refresh secret not configured")
	}

	token, err := jwt.ParseWithClaims(refreshToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshSecret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return "", errors.New("refresh token expired")
		}
	} else {
		return "", errors.New("invalid exp claim")
	}

	// Extract user ID from claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return "", errors.New("invalid user_id claim")
	}
	userID := uint(userIDFloat)

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("invalid username claim")
	}

	accessToken, _, err := jwtmw.GenerateToken(userID, username)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
