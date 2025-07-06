package postgres

import (
	"database/sql"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/model"
)

type MessageRepo struct {
	db *sql.DB
}

func NewMessageRepo(db *sql.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) Save(msg *model.Message) error {
	query := `
		INSERT INTO messages (chat_id, sender_id, text, ai_generated, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		msg.ChatID,
		msg.SenderID,
		msg.Content,
		msg.AIGenerated,
		time.Now(),
	).Scan(&msg.ID)

	return err
}

func (r *MessageRepo) GetByChat(chatID int64, limit, offset int) ([]model.Message, error) {
	query := `
		SELECT id, chat_id, sender_id, text, ai_generated, created_at
		FROM messages
		WHERE chat_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var msg model.Message
		err := rows.Scan(
			&msg.ID,
			&msg.ChatID,
			&msg.SenderID,
			&msg.Content,
			&msg.AIGenerated,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
