package ws

import (
	"context"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/proto/chatpb"
	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/ports"
)

func HandleAiAutoReply(
	msg kafkaadapter.AiAutoReplyPayload,
	chatClient chatpb.ChatServiceClient,
	hub ports.Hub,
) {
	resp, err := chatClient.GetUserWithChatByThreadID(context.Background(), &chatpb.GetUserWithChatByThreadIDRequest{
		ThreadId: msg.ThreadID,
	})
	if err != nil {
		log.Printf("failed to get user with chat by threadId %s: %v", msg.ThreadID, err)
		return
	}

	wsMsg := map[string]interface{}{
		"senderEmail": resp.UserEmail,
		"text":        msg.Text,
		"fromAI":      true,
	}

	hub.SendToRoom(resp.ChatId, "message", wsMsg)
}
