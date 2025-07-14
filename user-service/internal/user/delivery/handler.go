package grpc

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/proto/userpb"
	"github.com/Vovarama1992/go-ai-messenger/user-service/internal/user/ports"
)

type Handler struct {
	userpb.UnimplementedUserServiceServer
	service ports.UserService
}

func NewHandler(service ports.UserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.GetUserByEmailResponse, error) {
	user, err := h.service.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserByEmailResponse{
		Id:           user.ID,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (h *Handler) GetUserByID(ctx context.Context, req *userpb.GetUserByIDRequest) (*userpb.GetUserByIDResponse, error) {
	user, err := h.service.FindByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserByIDResponse{
		Id:    user.ID,
		Email: user.Email,
	}, nil
}

func (h *Handler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	user, err := h.service.Create(ctx, req.GetEmail(), req.GetPasswordHash())
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{
		Id: user.ID,
	}, nil
}
