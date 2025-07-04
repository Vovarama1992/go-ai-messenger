package httpadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/auth-service/internal/ports"
)

type Handler struct {
	service ports.AuthService
}

func NewHandler(service ports.AuthService) *Handler {
	return &Handler{service: service}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	UserID int64  `json:"user_id,omitempty"`
	Token  string `json:"access_token,omitempty"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	token, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokenResponse{Token: token})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := h.service.Register(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, "conflict", http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(tokenResponse{UserID: id})
}
