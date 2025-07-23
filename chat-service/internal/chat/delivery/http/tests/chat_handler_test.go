package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	handler "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/adapters/http"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/mocks"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/usecase"
	middleware "github.com/Vovarama1992/go-ai-messenger/pkg/authmiddleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateChat_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	mockBindingRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockAdvice := mocks.NewMockAdvicePublisher(ctrl)
	mockBroker := mocks.NewMockMessageBroker(ctrl)

	chatService := usecase.NewChatService(mockBroker, mockRepo, mockBindingRepo, mockAdvice)
	handler := handler.NewChatHandler(chatService, "invite-topic")

	router := chi.NewRouter()
	router.Post("/chats", handler.CreateChat)

	withMockAuth := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), middleware.UserIDKey, int64(123))
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	server := httptest.NewServer(withMockAuth(router))
	defer server.Close()

	// подготовим запрос
	body := `{"type":"private","members":[456]}`
	expectedChat := &model.Chat{
		ID:        1,
		CreatorID: 123,
		ChatType:  model.ChatTypePrivate,
		CreatedAt: time.Now().Unix(),
	}

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any(), []int64{456, 123}).
		DoAndReturn(func(ctx context.Context, chat *model.Chat, memberIDs []int64) error {
			chat.ID = expectedChat.ID
			chat.CreatorID = expectedChat.CreatorID
			chat.ChatType = expectedChat.ChatType
			chat.CreatedAt = expectedChat.CreatedAt
			return nil
		})

	resp, err := http.Post(server.URL+"/chats", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var got model.Chat
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.Equal(t, expectedChat.ID, got.ID)
	require.Equal(t, expectedChat.CreatorID, got.CreatorID)
	require.Equal(t, expectedChat.ChatType, got.ChatType)
}
