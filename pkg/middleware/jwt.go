package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// JWTClaims defines custom claims structure
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT access token and refresh token pair.
func GenerateToken(userID uint, username string) (accessToken string, refreshToken string, err error) {
	secret := []byte(os.Getenv("ACCESS_SECRET"))
	refreshSecret := []byte(os.Getenv("REFRESH_SECRET"))
	if len(secret) == 0 || len(refreshSecret) == 0 {
		return "", "", ErrMissingSecrets
	}

	// Access token
	atClaims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "access_token",
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err = at.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	// Refresh token
	rtClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   "refresh_token",
		ID:        username,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err = rt.SignedString(refreshSecret)
	return
}

// GenerateTokenFromRefresh validates a refresh token and issues new access and refresh tokens.
func GenerateTokenFromRefresh(refreshToken string) (string, string, error) {
	refreshSecret := []byte(os.Getenv("REFRESH_SECRET"))
	if len(refreshSecret) == 0 {
		return "", "", ErrMissingSecrets
	}

	// Parse refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil || !token.Valid {
		return "", "", err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || claims.Subject != "refresh_token" {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	// Extract userID from IssuedAt timestamp encoded as Unix seconds in Token ID if that's intended.
	// Here, assuming the username is stored in ID, and userID must be reloaded via service in real code.
	username := claims.ID

	// Generate new tokens
	return GenerateToken(parseUint(username), username)
}

// parseUint is a helper to convert string to uint
func parseUint(s string) uint {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(v)
}

// ErrMissingSecrets is returned when env secrets are not set
var ErrMissingSecrets = gin.Error{
	Err:  jwt.ErrTokenMalformed,
	Type: gin.ErrorTypePrivate,
}

// JWTAuthMiddleware validates the JWT access token.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := []byte(os.Getenv("ACCESS_SECRET"))
		if len(secret) == 0 {
			logrus.Error("ACCESS_SECRET env var not set")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			logrus.WithError(err).Warn("Invalid JWT token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}
		// set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
