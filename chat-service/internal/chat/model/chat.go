package model

import "errors"

type ChatType string

const (
	ChatTypePrivate ChatType = "private"
	ChatTypeGroup   ChatType = "group"
)

func (t ChatType) IsValid() error {
	switch t {
	case ChatTypePrivate, ChatTypeGroup:
		return nil
	default:
		return errors.New("invalid chat type")
	}
}

type Chat struct {
	ID        int64
	CreatorID int64
	Type      ChatType
	CreatedAt int64
}

type ChatMember struct {
	ChatID   int64
	UserID   int64
	Accepted bool
}
