package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type TelegramService struct {
	token  string
	chatID string
	client *http.Client
}

func NewTelegramService(token, chatID string) *TelegramService {
	return &TelegramService{
		token:  token,
		chatID: chatID,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *TelegramService) SendMessage(message string) error {
	if s.token == "" || s.chatID == "" {
		return fmt.Errorf("telegram token or chat ID is not configured")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.token)
	
	payload := map[string]interface{}{
		"chat_id":    s.chatID,
		"text":       message,
		"parse_mode": "HTML",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal telegram payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create telegram request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from telegram: %d", resp.StatusCode)
	}

	log.Printf("[INFO] Telegram message sent successfully")
	return nil
}
