package http

import (
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
)

type createChatRequest struct {
	Type    model.ChatType `json:"type" validate:"required"`
	Members []int64        `json:"members" validate:"omitempty,dive,gt=0"`
}

type sendInviteRequest struct {
	UserIDs []int64 `json:"userIds" validate:"required,dive,gt=0"`
}

type bindRequest struct {
	Type model.AIBindingType `json:"type" validate:"required"`
}
