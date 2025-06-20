package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/usecase"
	middleware "github.com/Vovarama1992/go-ai-messenger/pkg/authmiddleware"
)

type ChatHandler struct {
	service *usecase.ChatService
}

func NewChatHandler(service *usecase.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

type createChatRequest struct {
	Type model.ChatType `json:"type"`
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req createChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := req.Type.IsValid(); err != nil {
		http.Error(w, "invalid chat type", http.StatusBadRequest)
		return
	}

	chat, err := h.service.CreateChat(r.Context(), userID, req.Type)
	if err != nil {
		http.Error(w, "failed to create chat", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(chat)
}

func (h *ChatHandler) GetChatByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	chat, err := h.service.GetChatByID(r.Context(), id)
	if err != nil {
		http.Error(w, "chat not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(chat)
}

func (h *ChatHandler) RequestAdvice(w http.ResponseWriter, r *http.Request) {
	chatIDStr := chi.URLParam(r, "id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid chat ID", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.service.RequestAdvice(r.Context(), userID, chatID)
	if err != nil {
		http.Error(w, "failed to request advice: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
