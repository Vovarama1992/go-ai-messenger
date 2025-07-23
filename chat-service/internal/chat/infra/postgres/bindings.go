package postgres

import (
	"context"
	"database/sql"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/sony/gobreaker"
)

type ChatBindingRepo struct {
	db      *sql.DB
	breaker *gobreaker.CircuitBreaker
}

func NewChatBindingRepo(db *sql.DB, breaker *gobreaker.CircuitBreaker) *ChatBindingRepo {
	return &ChatBindingRepo{db: db, breaker: breaker}
}

func (r *ChatBindingRepo) Create(ctx context.Context, binding *model.ChatBinding) error {
	_, err := r.breaker.Execute(func() (interface{}, error) {
		query := `INSERT INTO chat_bindings (chat_id, user_id, type, created_at) VALUES ($1, $2, $3, $4)`
		_, err := r.db.ExecContext(ctx, query, binding.ChatID, binding.UserID, binding.BindingType, binding.CreatedAt)
		return nil, err
	})
	return err
}

func (r *ChatBindingRepo) FindByUserAndChat(ctx context.Context, userID, chatID int64) (*model.ChatBinding, error) {
	result, err := r.breaker.Execute(func() (interface{}, error) {
		query := `SELECT chat_id, user_id, type, created_at FROM chat_bindings WHERE chat_id = $1 AND user_id = $2`

		var binding model.ChatBinding
		err := r.db.QueryRowContext(ctx, query, chatID, userID).Scan(
			&binding.ChatID,
			&binding.UserID,
			&binding.BindingType,
			&binding.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		return &binding, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.ChatBinding), nil
}

func (r *ChatBindingRepo) Update(ctx context.Context, binding *model.ChatBinding) error {
	_, err := r.breaker.Execute(func() (interface{}, error) {
		query := `UPDATE chat_bindings SET type = $1, created_at = $2 WHERE chat_id = $3 AND user_id = $4`
		_, err := r.db.ExecContext(ctx, query, binding.BindingType, binding.CreatedAt, binding.ChatID, binding.UserID)
		return nil, err
	})
	return err
}

func (r *ChatBindingRepo) Delete(ctx context.Context, userID, chatID int64) error {
	_, err := r.breaker.Execute(func() (interface{}, error) {
		query := `DELETE FROM chat_bindings WHERE chat_id = $1 AND user_id = $2`
		_, err := r.db.ExecContext(ctx, query, chatID, userID)
		return nil, err
	})
	return err
}

func (r *ChatBindingRepo) FindBindingsByChatID(ctx context.Context, chatID int64) ([]*model.ChatBinding, error) {
	result, err := r.breaker.Execute(func() (interface{}, error) {
		query := `SELECT chat_id, user_id, type, created_at FROM chat_bindings WHERE chat_id = $1`

		rows, err := r.db.QueryContext(ctx, query, chatID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var bindings []*model.ChatBinding
		for rows.Next() {
			var b model.ChatBinding
			if err := rows.Scan(&b.ChatID, &b.UserID, &b.BindingType, &b.CreatedAt); err != nil {
				return nil, err
			}
			bindings = append(bindings, &b)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		return bindings, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*model.ChatBinding), nil
}

func (r *ChatBindingRepo) UpdateThreadID(ctx context.Context, chatID, userID int64, threadID string) error {
	_, err := r.breaker.Execute(func() (interface{}, error) {
		_, err := r.db.ExecContext(ctx, `
			UPDATE chat_bindings 
			SET thread_id = $1 
			WHERE chat_id = $2 AND user_id = $3
		`, threadID, chatID, userID)
		return nil, err
	})
	return err
}

func (r *ChatBindingRepo) FindByThreadID(ctx context.Context, threadID string) (*model.ChatBinding, error) {
	result, err := r.breaker.Execute(func() (interface{}, error) {
		query := `SELECT chat_id, user_id FROM chat_bindings WHERE thread_id = $1`

		var b model.ChatBinding
		err := r.db.QueryRowContext(ctx, query, threadID).Scan(
			&b.ChatID,
			&b.UserID,
		)
		if err != nil {
			return nil, err
		}

		return &b, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.ChatBinding), nil
}
