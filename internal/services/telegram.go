package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// TelegramMessage represents the Telegram sendMessage request body
type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

// SendNotification sends a message to Telegram
func SendNotification(message string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatID == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID not configured")
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	requestBody := TelegramMessage{
		ChatID: chatID,
		Text:   message,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Telegram API returned status: %d", resp.StatusCode)
	}

	return nil
}
