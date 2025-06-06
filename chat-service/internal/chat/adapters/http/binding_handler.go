package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/usecase"
	middleware "github.com/Vovarama1992/go-ai-messenger/pkg/authmiddleware"
)

type ChatBindingHandler struct {
	service *usecase.ChatBindingService
}

func NewChatBindingHandler(service *usecase.ChatBindingService) *ChatBindingHandler {
	return &ChatBindingHandler{service: service}
}

type bindRequest struct {
	Type model.AIBindingType `json:"type"`
}

func (h *ChatBindingHandler) CreateOrUpdateBinding(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.URL.Query().Get("chat_id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userEmail, ok := middleware.GetUserEmail(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req bindRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Type.IsValid(); err != nil {
		http.Error(w, "invalid binding type", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateBinding(r.Context(), userEmail, userID, chatID, req.Type)
	if err != nil {
		http.Error(w, "failed to bind AI", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ChatBindingHandler) DeleteBinding(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.URL.Query().Get("chat_id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.service.DeleteBinding(r.Context(), userID, chatID)
	if err != nil {
		http.Error(w, "failed to delete binding", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ChatBindingHandler) GetBinding(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.URL.Query().Get("chat_id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	binding, err := h.service.GetBinding(r.Context(), userID, chatID)
	if err != nil {
		http.Error(w, "binding not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(binding)
}
