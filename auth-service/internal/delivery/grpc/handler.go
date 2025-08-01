package grpc

import (
	"context"
	"errors"

	"github.com/Vovarama1992/go-ai-messenger/auth-service/internal/ports"
	"github.com/Vovarama1992/go-ai-messenger/proto/authpb"
)

type Handler struct {
	authpb.UnimplementedAuthServiceServer
	authService ports.AuthService
}

func NewHandler(authService ports.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	userID, email, err := h.authService.ValidateToken(ctx, req.GetToken())
	if err != nil {
		return nil, errors.New("invalid token")
	}

	return &authpb.ValidateTokenResponse{
		UserId: userID,
		Email:  email,
	}, nil
}
