package postgres

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/user-service/internal/user/model"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRepository struct {
	db      *pgxpool.Pool
	breaker *gobreaker.CircuitBreaker
}

func NewUserRepository(db *pgxpool.Pool, breaker *gobreaker.CircuitBreaker) *UserRepository {
	return &UserRepository{
		db:      db,
		breaker: breaker,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	result, err := r.breaker.Execute(func() (interface{}, error) {
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
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.User), nil
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) (*model.User, error) {
	if len(email) > 320 {
		return nil, status.Error(codes.InvalidArgument, "email too long")
	}
	result, err := r.breaker.Execute(func() (interface{}, error) {
		const query = `
			INSERT INTO users (email, password_hash)
			VALUES ($1, $2)
			RETURNING id
		`
		var id int64
		err := r.db.QueryRow(ctx, query, email, passwordHash).Scan(&id)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
				return nil, status.Error(codes.AlreadyExists, "email already registered")
			}
			return nil, err
		}

		return &model.User{
			ID:           id,
			Email:        email,
			PasswordHash: passwordHash,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.User), nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	result, err := r.breaker.Execute(func() (interface{}, error) {
		const query = `
			SELECT id, email, password_hash
			FROM users
			WHERE id = $1
		`
		row := r.db.QueryRow(ctx, query, id)
		var user model.User
		err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
		if err != nil {
			return nil, err
		}
		return &user, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.User), nil
}
