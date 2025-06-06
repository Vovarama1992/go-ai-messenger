package grpc

import (
	"context"

	usecase "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/usecase"
	messagepb "github.com/Vovarama1992/go-ai-messenger/proto/messagepb"
)

type MessageHandler struct {
	messagepb.UnimplementedMessageServiceServer
	messageService *usecase.MessageService
}

func NewMessageHandler(messageService *usecase.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) GetMessagesByChat(ctx context.Context, req *messagepb.GetMessagesRequest) (*messagepb.GetMessagesResponse, error) {
	messages, err := h.messageService.GetMessagesByChat(ctx, req.ChatId)
	if err != nil {
		return nil, err
	}

	return &messagepb.GetMessagesResponse{
		Messages: messages,
	}, nil
}
