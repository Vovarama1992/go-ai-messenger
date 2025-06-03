package grpc

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/usecase"
	"github.com/Vovarama1992/go-ai-messenger/proto/chatpb"
)

type ChatHandler struct {
	chatpb.UnimplementedChatServiceServer
	chatService    *usecase.ChatService
	bindingService *usecase.ChatBindingService
}

func NewChatHandler(chatService *usecase.ChatService, bindingService *usecase.ChatBindingService) *ChatHandler {
	return &ChatHandler{
		chatService:    chatService,
		bindingService: bindingService,
	}
}

func (h *ChatHandler) GetChatByID(ctx context.Context, req *chatpb.GetChatByIDRequest) (*chatpb.GetChatByIDResponse, error) {
	chat, err := h.chatService.GetChatByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &chatpb.GetChatByIDResponse{
		Id:        chat.ID,
		CreatorId: chat.CreatorID,
		Type:      string(chat.Type),
		CreatedAt: chat.CreatedAt,
	}, nil
}

func (h *ChatHandler) GetBindingsByChat(ctx context.Context, req *chatpb.GetBindingsByChatRequest) (*chatpb.GetBindingsByChatResponse, error) {
	bindings, err := h.bindingService.GetBindingsByChat(ctx, req.ChatId)
	if err != nil {
		return nil, err
	}

	var resp chatpb.GetBindingsByChatResponse
	for _, b := range bindings {
		resp.Bindings = append(resp.Bindings, &chatpb.ChatBinding{
			UserId: b.UserID,
			Type:   string(b.Type),
		})
	}

	return &resp, nil
}
