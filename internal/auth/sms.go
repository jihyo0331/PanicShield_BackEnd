package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type otpEntry struct {
	code      string
	expiresAt time.Time
}

var (
	otpStore      = make(map[uint]otpEntry)
	otpStoreMutex sync.RWMutex
)

func init() {
	rand.Seed(time.Now().UnixNano())
	// Optional: background cleanup goroutine
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			now := time.Now()
			otpStoreMutex.Lock()
			for userID, entry := range otpStore {
				if entry.expiresAt.Before(now) {
					delete(otpStore, userID)
				}
			}
			otpStoreMutex.Unlock()
		}
	}()
}

// GenerateAndSendOTP generates a 6-digit code, stores it, and sends it via SMS.
func GenerateAndSendOTP(userID uint, phone string) error {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiry := time.Now().Add(5 * time.Minute)

	otpStoreMutex.Lock()
	otpStore[userID] = otpEntry{code: code, expiresAt: expiry}
	otpStoreMutex.Unlock()

	payload := map[string]string{
		"to":      phone,
		"message": fmt.Sprintf("Your verification code is %s", code),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal SMS payload")
		return err
	}

	apiURL := os.Getenv("SMS_API_URL")
	apiKey := os.Getenv("SMS_API_KEY")
	if apiURL == "" || apiKey == "" {
		logrus.Error("Missing SMS_API_URL or SMS_API_KEY environment variable")
		return errors.New("missing SMS API configuration")
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		logrus.WithError(err).Error("Failed to create SMS API request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to send SMS OTP")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"phone":  phone,
		}).Info("OTP sent successfully")
		return nil
	}
	logrus.WithFields(logrus.Fields{
		"userID":     userID,
		"phone":      phone,
		"statusCode": resp.StatusCode,
	}).Error("Failed to send OTP SMS")
	return fmt.Errorf("failed to send OTP SMS, status: %s", resp.Status)
}

// ValidateOTP checks if the code for userID is valid and removes it on success.
func ValidateOTP(userID uint, code string) bool {
	otpStoreMutex.Lock()
	defer otpStoreMutex.Unlock()
	entry, exists := otpStore[userID]
	if !exists {
		return false
	}
	if time.Now().After(entry.expiresAt) {
		delete(otpStore, userID)
		return false
	}
	if entry.code != code {
		return false
	}
	delete(otpStore, userID)
	return true
}
