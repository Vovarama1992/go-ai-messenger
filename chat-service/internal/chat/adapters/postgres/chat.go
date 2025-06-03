package postgres

import (
	"context"
	"database/sql"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
)

type ChatRepo struct {
	db *sql.DB
}

func NewChatRepo(db *sql.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

func (r *ChatRepo) Create(ctx context.Context, chat *model.Chat) error {
	query := `INSERT INTO chats (creator_id, type, created_at) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, chat.CreatorID, chat.Type, chat.CreatedAt).Scan(&chat.ID)
}

func (r *ChatRepo) FindByID(ctx context.Context, id int64) (*model.Chat, error) {
	query := `SELECT id, creator_id, type, created_at FROM chats WHERE id = $1`

	var chat model.Chat
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&chat.ID,
		&chat.CreatorID,
		&chat.Type,
		&chat.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}
