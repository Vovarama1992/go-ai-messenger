package http_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/delivery/http"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/mocks"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/usecase"
	middleware "github.com/Vovarama1992/go-ai-messenger/pkg/authmiddleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateOrUpdateBinding_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockBroker := mocks.NewMockMessageBroker(ctrl)
	mockClient := mocks.NewMockMessageClient(ctrl)

	service := usecase.NewChatBindingService(mockRepo, mockBroker, mockClient)
	h := handler.NewChatBindingHandler(service)

	router := chi.NewRouter()
	router.Post("/bind", h.CreateOrUpdateBinding)

	withMockAuth := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), middleware.UserIDKey, int64(123))
			ctx = context.WithValue(ctx, middleware.EmailKey, "user@example.com")
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	server := httptest.NewServer(withMockAuth(router))
	defer server.Close()

	chatID := int64(42)
	body := []byte(`{"type":"advice"}`)

	// симулируем: привязки нет, GetMessagesByChat вернёт пусто
	mockRepo.EXPECT().
		FindByUserAndChat(gomock.Any(), int64(123), chatID).
		Return(nil, context.DeadlineExceeded) // чтобы пошёл через Create

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	mockClient.EXPECT().
		GetMessagesByChat(gomock.Any(), chatID).
		Return(nil, nil)

	mockBroker.EXPECT().
		SendAiBindingInit(gomock.Any(), gomock.Any()).
		Return(nil)

	req, err := http.NewRequest("POST", server.URL+"/bind?chat_id=42", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}
