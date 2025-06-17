package chatbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"ps_backend/model"
)

type ChatbotService struct {
	db *gorm.DB
}

func NewChatbotService(db *gorm.DB) *ChatbotService {
	return &ChatbotService{db: db}
}

func (s *ChatbotService) SendMessage(userID uint, message string) (string, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		logrus.Errorf("failed to load user %d: %v", userID, err)
		return "", fmt.Errorf("user not found")
	}

	var history []model.ChatbotLog
	if err := s.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(5).
		Find(&history).Error; err != nil {
		logrus.Errorf("failed to load chat history for user %d: %v", userID, err)
		return "", fmt.Errorf("failed to load chat history")
	}

	// Reverse history to chronological order
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	promptBuilder := &bytes.Buffer{}
	fmt.Fprintf(promptBuilder, "You are chatting with user %s.\n", user.Username)
	fmt.Fprintf(promptBuilder, "Style: %s\nTone: %s\n\n", user.SpeakingStyle, user.Tone)
	fmt.Fprintf(promptBuilder, "Conversation history:\n")
	for _, log := range history {
		role := "User"
		if log.Sender == "bot" {
			role = "Bot"
		}
		fmt.Fprintf(promptBuilder, "%s: %s\n", role, log.Message)
	}
	fmt.Fprintf(promptBuilder, "User: %s\nBot:", message)
	prompt := promptBuilder.String()

	apiURL := os.Getenv("GEMINI_API_URL")
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiURL == "" || apiKey == "" {
		logrus.Error("GEMINI_API_URL or GEMINI_API_KEY environment variable is not set")
		return "", fmt.Errorf("chatbot service is not properly configured")
	}

	type geminiRequest struct {
		Model       string  `json:"model"`
		Prompt      string  `json:"prompt"`
		MaxTokens   int     `json:"max_tokens"`
		Temperature float64 `json:"temperature"`
	}

	reqBody := geminiRequest{
		Model:       "gemini-1",
		Prompt:      prompt,
		MaxTokens:   512,
		Temperature: 0.7,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		logrus.Errorf("failed to marshal Gemini request: %v", err)
		return "", fmt.Errorf("internal error")
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		logrus.Errorf("failed to create Gemini request: %v", err)
		return "", fmt.Errorf("internal error")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("failed to call Gemini API: %v", err)
		return "", fmt.Errorf("failed to get response from chatbot")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Gemini API returned non-200 status: %d", resp.StatusCode)
		return "", fmt.Errorf("chatbot service error")
	}

	type geminiResponse struct {
		Reply string `json:"reply"`
	}

	var respData geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		logrus.Errorf("failed to decode Gemini response: %v", err)
		return "", fmt.Errorf("invalid response from chatbot")
	}

	now := time.Now()
	userLog := model.ChatbotLog{
		UserID:    userID,
		Message:   message,
		Sender:    "user",
		CreatedAt: now,
	}
	botLog := model.ChatbotLog{
		UserID:    userID,
		Message:   respData.Reply,
		Sender:    "bot",
		CreatedAt: now,
	}

	if err := s.db.Create(&userLog).Error; err != nil {
		logrus.Errorf("failed to save user message log: %v", err)
		return "", fmt.Errorf("internal error")
	}
	if err := s.db.Create(&botLog).Error; err != nil {
		logrus.Errorf("failed to save bot reply log: %v", err)
		return "", fmt.Errorf("internal error")
	}

	logrus.Infof("chatbot reply sent to user %d", userID)
	return respData.Reply, nil
}
