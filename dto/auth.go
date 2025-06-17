package dto

// SignUpRequest represents the JSON body for user registration.
type SignUpRequest struct {
	Username      string `json:"username" binding:"required,min=2,max=32"`
	Password      string `json:"password" binding:"required,min=6"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
	SpeakingStyle string `json:"speaking_style" binding:"required"`
	Tone          string `json:"tone" binding:"required"`
}

// LoginRequest represents the JSON body for user login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// VerifyPhoneRequest represents the JSON body for phone verification.
type VerifyPhoneRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Code   string `json:"code" binding:"required,len=6"`
}

// RefreshTokenRequest represents the JSON body to refresh an auth token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
