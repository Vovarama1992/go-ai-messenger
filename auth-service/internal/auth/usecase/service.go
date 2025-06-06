package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/auth-service/internal/auth/ports"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userClient ports.UserClient
	jwtSecret  string
}

func NewAuthService(userClient ports.UserClient, jwtSecret string) *AuthServiceImpl {
	return &AuthServiceImpl{
		userClient: userClient,
		jwtSecret:  jwtSecret,
	}
}

func (s *AuthServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	id, passwordHash, err := s.userClient.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub":   id,
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthServiceImpl) Register(ctx context.Context, email, password string) (int64, error) {
	_, _, err := s.userClient.GetByEmail(ctx, email)
	if err == nil {
		return 0, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	return s.userClient.Create(ctx, email, string(hash))
}

func (s *AuthServiceImpl) ValidateToken(ctx context.Context, tokenString string) (int64, string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid claims")
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, "", errors.New("sub claim missing or invalid")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return 0, "", errors.New("email claim missing or invalid")
	}

	return int64(sub), email, nil
}
