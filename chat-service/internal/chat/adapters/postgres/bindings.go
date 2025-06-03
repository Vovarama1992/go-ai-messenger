package postgres

import (
	"context"
	"database/sql"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
)

type ChatBindingRepo struct {
	db *sql.DB
}

func NewChatBindingRepo(db *sql.DB) *ChatBindingRepo {
	return &ChatBindingRepo{db: db}
}

func (r *ChatBindingRepo) Create(ctx context.Context, binding *model.ChatBinding) error {
	query := `INSERT INTO chat_bindings (chat_id, user_id, type, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, binding.ChatID, binding.UserID, binding.Type, binding.CreatedAt)
	return err
}

func (r *ChatBindingRepo) FindByUserAndChat(ctx context.Context, userID, chatID int64) (*model.ChatBinding, error) {
	query := `SELECT chat_id, user_id, type, created_at FROM chat_bindings WHERE chat_id = $1 AND user_id = $2`

	var binding model.ChatBinding
	err := r.db.QueryRowContext(ctx, query, chatID, userID).Scan(
		&binding.ChatID,
		&binding.UserID,
		&binding.Type,
		&binding.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &binding, nil
}

func (r *ChatBindingRepo) Update(ctx context.Context, binding *model.ChatBinding) error {
	query := `UPDATE chat_bindings SET type = $1, created_at = $2 WHERE chat_id = $3 AND user_id = $4`
	_, err := r.db.ExecContext(ctx, query, binding.Type, binding.CreatedAt, binding.ChatID, binding.UserID)
	return err
}

func (r *ChatBindingRepo) Delete(ctx context.Context, userID, chatID int64) error {
	query := `DELETE FROM chat_bindings WHERE chat_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, chatID, userID)
	return err
}

func (r *ChatBindingRepo) FindBindingsByChatID(ctx context.Context, chatID int64) ([]*model.ChatBinding, error) {
	query := `SELECT chat_id, user_id, type, created_at FROM chat_bindings WHERE chat_id = $1`

	rows, err := r.db.QueryContext(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bindings []*model.ChatBinding
	for rows.Next() {
		var b model.ChatBinding
		if err := rows.Scan(&b.ChatID, &b.UserID, &b.Type, &b.CreatedAt); err != nil {
			return nil, err
		}
		bindings = append(bindings, &b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bindings, nil
}
