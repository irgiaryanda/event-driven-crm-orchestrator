package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// LLMRequest represents the OpenAI-compatible request body
type LLMRequest struct {
	Model    string       `json:"model"`
	Messages []LLMMessage `json:"messages"`
}

// LLMMessage represents a message in the LLM request
type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMResponse represents the OpenAI-compatible response body
type LLMResponse struct {
	Choices []LLMChoice `json:"choices"`
}

// LLMChoice represents a choice in the LLM response
type LLMChoice struct {
	Message LLMMessage `json:"message"`
}

// CategorizePayload sends payload to LLM and returns the category
func CategorizePayload(payload string) (string, error) {
	apiURL := os.Getenv("LLM_API_URL")
	apiKey := os.Getenv("LLM_API_KEY")
	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = "gpt-4o"
	}

	if apiURL == "" || apiKey == "" {
		return "", fmt.Errorf("LLM_API_URL or LLM_API_KEY not configured")
	}

	systemPrompt := "You are an intelligent CRM router. Categorize the following JSON payload into one of these exact categories: SUPPORT, SALES, BILLING, or OTHER. The payload may be written in English, Indonesian, or any other language. Understand the context regardless of the language and reply strictly with the category name only (e.g., SUPPORT). Do not include any other words."

	requestBody := LLMRequest{
		Model: model,
		Messages: []LLMMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: payload},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM API returned status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var llmResp LLMResponse
	if err := json.NewDecoder(resp.Body).Decode(&llmResp); err != nil {
		return "", err
	}

	if len(llmResp.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return llmResp.Choices[0].Message.Content, nil
}
