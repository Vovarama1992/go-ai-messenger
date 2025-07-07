package ws_test

import (
	"testing"

	"github.com/Vovarama1992/go-ai-messenger/proto/chatpb"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/delivery/ws"
	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestHandleAiAutoReply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatClient := mocks.NewMockChatServiceClient(ctrl)
	mockHub := mocks.NewMockHub(ctrl)

	testThreadID := "thread-123"
	testText := "hello from AI"
	testChatID := int64(42)
	testUserEmail := "user@example.com"

	// Ожидаем вызов GetUserWithChatByThreadID и возвращаем ответ
	mockChatClient.EXPECT().
		GetUserWithChatByThreadID(gomock.Any(), &chatpb.GetUserWithChatByThreadIDRequest{ThreadId: testThreadID}).
		Return(&chatpb.GetUserWithChatByThreadIDResponse{
			ChatId:    testChatID,
			UserEmail: testUserEmail,
		}, nil).
		Times(1)

	// Ожидаем, что hub.SendToRoom будет вызван с правильными параметрами
	mockHub.EXPECT().
		SendToRoom(testChatID, "message", gomock.Any()).
		DoAndReturn(func(chatID int64, event string, data any) {
			if chatID != testChatID {
				t.Errorf("unexpected chatID: got %v want %v", chatID, testChatID)
			}
			if event != "message" {
				t.Errorf("unexpected event: got %v want 'message'", event)
			}
			// data должен содержать senderEmail, text и fromAI=true
			msgMap, ok := data.(map[string]interface{})
			if !ok {
				t.Errorf("data is not a map[string]interface{}")
			}
			if msgMap["senderEmail"] != testUserEmail {
				t.Errorf("unexpected senderEmail: got %v want %v", msgMap["senderEmail"], testUserEmail)
			}
			if msgMap["text"] != testText {
				t.Errorf("unexpected text: got %v want %v", msgMap["text"], testText)
			}
			if msgMap["fromAI"] != true {
				t.Errorf("expected fromAI to be true")
			}
		}).
		Times(1)

	// Вызываем тестируемую функцию
	ws.HandleAiAutoReply(
		kafkaadapter.AiAutoReplyPayload{
			ThreadID: testThreadID,
			Text:     testText,
		},
		mockChatClient,
		mockHub,
	)
}
