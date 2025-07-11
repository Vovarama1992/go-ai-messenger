package gptutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetLastAssistantMessage(ctx context.Context, apiKey, baseURL, threadID string) (string, error) {
	url := fmt.Sprintf("%s/v1/threads/%s/messages?order=desc&limit=1", baseURL, threadID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch messages manually: %w", err)
	}
	defer resp.Body.Close()

	var data struct {
		Messages []struct {
			Role    string `json:"role"`
			Content []struct {
				Text struct {
					Value string `json:"value"`
				} `json:"text"`
			} `json:"content"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode assistant response: %w", err)
	}

	if len(data.Messages) == 0 || data.Messages[0].Role != "assistant" {
		return "", fmt.Errorf("no assistant response found")
	}

	return data.Messages[0].Content[0].Text.Value, nil
}
