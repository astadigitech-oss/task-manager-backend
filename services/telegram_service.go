package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TelegramService interface {
	SendNotification(chatID string, message string) error
}

type telegramService struct {
	botToken string
	client   *http.Client
}

func NewTelegramService(botToken string) TelegramService {
	return &telegramService{
		botToken: botToken,
		client:   &http.Client{},
	}
}

func (s *telegramService) SendNotification(chatID string, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)

	requestBody, err := json.Marshal(map[string]string{
		"chat_id": chatID,
		"text":    message,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification, status code: %d", resp.StatusCode)
	}

	return nil
}
