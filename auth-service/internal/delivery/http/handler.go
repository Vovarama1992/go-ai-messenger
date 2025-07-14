package httpadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/auth-service/internal/ports"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	service ports.AuthService
}

func NewHandler(service ports.AuthService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	token, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(TokenResponse{Email: req.Email, Token: token})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := h.service.Register(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, "conflict", http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(TokenResponse{UserID: id, Email: req.Email})
}
