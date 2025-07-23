package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"ps_backend/db"
	"ps_backend/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var chatdb = db.GetDB()

type GeminiRequest struct {
	Contents []struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type ChatRequest struct {
	UserID  uint   `json:"user_id"`
	Message string `json:"message"`
}

func ChatWithGemini(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "데이터 형식 오류"})
		return
	}

	// 유저 말투, 스타일 로드
	var user model.User
	if err := chatdb.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "유저 없음"})
		return
	}

	// 최근 대화 이력(최대 5개)
	var logs []model.ChatbotLog
	chatdb.Where("user_id = ?", req.UserID).Order("created_at desc").Limit(5).Find(&logs)
	var history string
	for i := len(logs) - 1; i >= 0; i-- {
		history += fmt.Sprintf("%s: %s\n", logs[i].Sender, logs[i].Message)
	}

	// Gemini 프롬프트 구성
	prompt := fmt.Sprintf(
		"사용자의 말투는 '%s', 스타일은 '%s'입니다.\n최근 대화 기록:\n%sUser: %s",
		user.SpeakingStyle, user.Tone, history, req.Message,
	)

	// Gemini API 호출
	reply, err := CallGeminiAPI(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gemini 호출 실패"})
		return
	}

	// 로그 저장
	chatdb.Create(&model.ChatbotLog{UserID: req.UserID, Message: req.Message, Sender: "user"})
	chatdb.Create(&model.ChatbotLog{UserID: req.UserID, Message: reply, Sender: "bot"})

	c.JSON(http.StatusOK, gin.H{"reply": reply})
}

// GetChatHistory retrieves the chat history for a given user_id query parameter.
func GetChatHistory(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "유효하지 않은 user_id"})
		return
	}

	// Fetch all chat logs for the user, ordered by creation time ascending
	var logs []model.ChatbotLog
	if err := chatdb.Where("user_id = ?", userID).
		Order("created_at asc").
		Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "대화 기록 조회 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

func CallGeminiAPI(prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY 환경변수가 설정되지 않음")
	}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)
	payload := GeminiRequest{
		Contents: []struct {
			Role  string `json:"role"`
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Role: "user",
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}
	jsonBytes, _ := json.Marshal(payload)
	logrus.Infof("[Gemini Debug] Request URL: %s", url)
	logrus.Infof("[Gemini Debug] Request Payload: %s", string(jsonBytes))
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	logrus.Infof("[Gemini Debug] HTTP Request sent, awaiting response...")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	logrus.Infof("[Gemini Debug] Response Status: %s", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Infof("[Gemini Debug] Response Body: %s", string(body))
	// 응답에서 답변 텍스트만 추출
	var gemResp GeminiResponse
	if err := json.Unmarshal(body, &gemResp); err != nil {
		return "", err
	}
	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("Gemini 응답 파싱 실패")
	}
	return gemResp.Candidates[0].Content.Parts[0].Text, nil
}
