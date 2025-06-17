package dto

// VitalRequest represents the JSON body for registering a vital sign.
type VitalRequest struct {
	UserID      uint `json:"user_id" binding:"required"`
	HeartRate   int  `json:"heart_rate" binding:"required"`
	BreathRate  int  `json:"breath_rate" binding:"required"`
	StressLevel int  `json:"stress_level" binding:"required"`
}
