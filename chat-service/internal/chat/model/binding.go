package model

import "errors"

type AIBindingType string

const (
	AIBindingAdvice    AIBindingType = "advice"
	AIBindingAutoreply AIBindingType = "autoreply"
)

func (a AIBindingType) IsValid() error {
	switch a {
	case AIBindingAdvice, AIBindingAutoreply:
		return nil
	default:
		return errors.New("invalid AI binding type")
	}
}

type ChatBinding struct {
	ChatID    int64
	UserID    int64
	Type      AIBindingType
	CreatedAt int64
}
