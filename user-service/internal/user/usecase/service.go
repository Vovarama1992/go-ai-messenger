package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/Vovarama1992/go-ai-messenger/user-service/internal/user/model"
	"github.com/Vovarama1992/go-ai-messenger/user-service/internal/user/ports"
	"github.com/Vovarama1992/go-utils/ctxutil"
)

var (
	ErrEmailRequired    = errors.New("email обязателен")
	ErrPasswordRequired = errors.New("password обязателен")
	ErrEmailTaken       = errors.New("email уже используется")
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, 2)
	defer cancel()
	return s.repo.FindByEmail(ctx, email)
}

func (s *UserService) FindByID(ctx context.Context, id int64) (*model.User, error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, 2)
	defer cancel()
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, email, passwordHash string) (*model.User, error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, 2)
	defer cancel()
	email = strings.TrimSpace(email)
	passwordHash = strings.TrimSpace(passwordHash)

	if email == "" {
		return nil, ErrEmailRequired
	}
	if passwordHash == "" {
		return nil, ErrPasswordRequired
	}

	if existing, _ := s.repo.FindByEmail(ctx, email); existing != nil {
		return nil, ErrEmailTaken
	}

	return s.repo.Create(ctx, email, passwordHash)
}
