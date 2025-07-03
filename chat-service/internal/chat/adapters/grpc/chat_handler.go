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
	userService    *usecase.UserService
}

var _ chatpb.ChatServiceServer = (*ChatHandler)(nil)

func NewChatHandler(chatService *usecase.ChatService, bindingService *usecase.ChatBindingService, userService *usecase.UserService) *ChatHandler {
	return &ChatHandler{
		chatService:    chatService,
		bindingService: bindingService,
		userService:    userService,
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

func (h *ChatHandler) GetUserWithChatByThreadID(ctx context.Context, req *chatpb.GetUserWithChatByThreadIDRequest) (*chatpb.GetUserWithChatByThreadIDResponse, error) {
	binding, err := h.bindingService.FindByThreadID(ctx, req.ThreadId)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.GetUserByID(ctx, binding.UserID)
	if err != nil {
		return nil, err
	}

	return &chatpb.GetUserWithChatByThreadIDResponse{
		UserId:    binding.UserID,
		ChatId:    binding.ChatID,
		UserEmail: user.Email,
	}, nil
}

func (h *ChatHandler) GetUsersByChatID(ctx context.Context, req *chatpb.GetUsersByChatIDRequest) (*chatpb.GetUsersByChatIDResponse, error) {
	userIDs, err := h.chatService.GetParticipants(ctx, req.ChatId)
	if err != nil {
		return nil, err
	}

	return &chatpb.GetUsersByChatIDResponse{
		UserIds: userIDs,
	}, nil
}

func (h *ChatHandler) GetThreadContext(ctx context.Context, req *chatpb.GetThreadContextRequest) (*chatpb.GetThreadContextResponse, error) {
	// По threadID получаем биндинг
	binding, err := h.bindingService.FindByThreadID(ctx, req.ThreadId)
	if err != nil {
		return nil, err
	}

	// Получаем юзера с email по userID
	user, err := h.userService.GetUserByID(ctx, binding.UserID)
	if err != nil {
		return nil, err
	}

	// Получаем участников чата по chatID
	userIDs, err := h.chatService.GetParticipants(ctx, binding.ChatID)
	if err != nil {
		return nil, err
	}

	return &chatpb.GetThreadContextResponse{
		SenderUserId:    binding.UserID,
		SenderUserEmail: user.Email,
		ChatId:          binding.ChatID,
		ChatUserIds:     userIDs,
	}, nil
}
