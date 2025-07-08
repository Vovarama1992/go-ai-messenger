//go:build integration

package gpt_test

import (
	"context"
	"os"
	"testing"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/infra/gpt"
)

func TestSmoke_CreateThread(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	assistantID := os.Getenv("OPENAI_ASSISTANT_ID")

	if apiKey == "" || assistantID == "" {
		t.Skip("OPENAI_API_KEY or OPENAI_ASSISTANT_ID not set")
	}

	client := gpt.NewClient(apiKey, assistantID)

	threadID, err := client.CreateThreadForUserAndChat(
		context.Background(),
		"user@example.com",
		nil,
	)

	if err != nil {
		t.Fatalf("failed to create thread: %v", err)
	}

	if threadID == "" {
		t.Fatal("thread ID is empty")
	}
}
