package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/sony/gobreaker"
)

type ChatRepo struct {
	db      *sql.DB
	breaker *gobreaker.CircuitBreaker
}

func NewChatRepo(db *sql.DB, breaker *gobreaker.CircuitBreaker) *ChatRepo {
	return &ChatRepo{db: db, breaker: breaker}
}

func (r *ChatRepo) Create(ctx context.Context, chat *model.Chat, memberIDs []int64) error {
	_, err := r.breaker.Execute(func() (interface{}, error) {
		query := `INSERT INTO chats (creator_id, type, created_at) VALUES ($1, $2, $3) RETURNING id`
		err := r.db.QueryRowContext(ctx, query, chat.CreatorID, chat.ChatType, chat.CreatedAt).Scan(&chat.ID)
		if err != nil {
			return nil, err
		}

		if len(memberIDs) > 0 {
			values := ""
			args := []interface{}{}
			for i, userID := range memberIDs {
				values += fmt.Sprintf("($1, $%d, true),", i+2)
				args = append(args, userID)
			}
			query = "INSERT INTO chat_members (chat_id, user_id, accepted) VALUES " + values[:len(values)-1] + " ON CONFLICT DO NOTHING"
			args = append([]interface{}{chat.ID}, args...)
			_, err = r.db.ExecContext(ctx, query, args...)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})
	return err
}

func (r *ChatRepo) FindByID(ctx context.Context, id int64) (*model.Chat, error) {
	result, err := r.breaker.Execute(func() (interface{}, error) {
		query := `SELECT id, creator_id, type, created_at FROM chats WHERE id = $1`
		var chat model.Chat
		err := r.db.QueryRowContext(ctx, query, id).Scan(&chat.ID, &chat.CreatorID, &chat.ChatType, &chat.CreatedAt)
		if err != nil {
			return nil, err
		}
		return &chat, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.Chat), nil
}

func (r *ChatRepo) SendInvite(ctx context.Context, chatID int64, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}

	_, err := r.breaker.Execute(func() (interface{}, error) {
		query := `INSERT INTO chat_members (chat_id, user_id, accepted) VALUES `
		args := []interface{}{chatID}
		values := ""

		for i, userID := range userIDs {
			values += `($1, $` + strconv.Itoa(i+2) + `, false),`
			args = append(args, userID)
		}

		query += values[:len(values)-1] + ` ON CONFLICT (chat_id, user_id) DO NOTHING`

		_, err := r.db.ExecContext(ctx, query, args...)
		return nil, err
	})
	return err
}

func (r *ChatRepo) AcceptInvite(ctx context.Context, chatID, userID int64) error {
	_, err := r.breaker.Execute(func() (interface{}, error) {
		query := `UPDATE chat_members SET accepted = true WHERE chat_id = $1 AND user_id = $2`
		_, err := r.db.ExecContext(ctx, query, chatID, userID)
		return nil, err
	})
	return err
}

func (r *ChatRepo) GetChatParticipants(ctx context.Context, chatID int64) ([]int64, error) {
	result, err := r.breaker.Execute(func() (interface{}, error) {
		query := `SELECT user_id FROM chat_members WHERE chat_id = $1`
		rows, err := r.db.QueryContext(ctx, query, chatID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var participants []int64
		for rows.Next() {
			var userID int64
			if err := rows.Scan(&userID); err != nil {
				return nil, err
			}
			participants = append(participants, userID)
		}
		return participants, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]int64), nil
}

func (r *ChatRepo) GetPendingInvites(ctx context.Context, userID int64) ([]model.Chat, error) {
	result, err := r.breaker.Execute(func() (interface{}, error) {
		query := `
			SELECT c.id, c.creator_id, c.type, c.created_at
			FROM chats c
			INNER JOIN chat_members cm ON cm.chat_id = c.id
			WHERE cm.user_id = $1 AND cm.accepted = false
			ORDER BY c.created_at DESC
		`
		rows, err := r.db.QueryContext(ctx, query, userID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var chats []model.Chat
		for rows.Next() {
			var c model.Chat
			if err := rows.Scan(&c.ID, &c.CreatorID, &c.ChatType, &c.CreatedAt); err != nil {
				return nil, err
			}
			chats = append(chats, c)
		}
		return chats, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]model.Chat), nil
}
