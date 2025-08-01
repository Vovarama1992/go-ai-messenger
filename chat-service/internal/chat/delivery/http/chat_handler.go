package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/usecase"
	middleware "github.com/Vovarama1992/go-ai-messenger/pkg/authmiddleware"
)

type ChatHandler struct {
	service     *usecase.ChatService
	inviteTopic string
}

func NewChatHandler(service *usecase.ChatService, inviteTopic string) *ChatHandler {
	return &ChatHandler{service: service, inviteTopic: inviteTopic}
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
	if err := validate.Struct(req); err != nil {
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// добавим creator, если его нет в списке
	members := req.Members
	found := false
	for _, id := range members {
		if id == userID {
			found = true
			break
		}
	}
	if !found {
		members = append(members, userID)
	}

	chat, err := h.service.CreateChat(r.Context(), userID, req.Type, members)
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

func (h *ChatHandler) SendInvite(w http.ResponseWriter, r *http.Request) {
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

	participants, err := h.service.GetParticipants(r.Context(), chatID)
	if err != nil {
		http.Error(w, "failed to check chat participants", http.StatusInternalServerError)
		return
	}

	found := false
	for _, id := range participants {
		if id == userID {
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "forbidden: user is not participant of chat", http.StatusForbidden)
		return
	}

	var req sendInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.SendInvite(r.Context(), chatID, req.UserIDs, h.inviteTopic)
	if err != nil {
		http.Error(w, "failed to send invites: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ChatHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.AcceptInvite(r.Context(), chatID, userID)
	if err != nil {
		http.Error(w, "failed to accept invite: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ChatHandler) GetPendingInvites(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	chats, err := h.service.GetPendingInvites(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get pending invites: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(chats)
}
