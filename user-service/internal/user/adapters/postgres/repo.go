package postgres

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/user-service/internal/user/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	const query = `
		SELECT id, email, password_hash
		FROM users
		WHERE email = $1
	`

	row := r.db.QueryRow(ctx, query, email)

	var user model.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) (*model.User, error) {
	const query = `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRow(ctx, query, email, passwordHash).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}
