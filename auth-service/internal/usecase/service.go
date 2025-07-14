package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/auth-service/internal/ports"
	"github.com/Vovarama1992/go-utils/grpcutil"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceImpl struct {
	userClient    ports.UserClient
	jwtSecret     string
	bcryptLimiter chan struct{}
}

func NewAuthService(userClient ports.UserClient, jwtSecret string, bcryptLimiter chan struct{}) *AuthServiceImpl {
	return &AuthServiceImpl{
		userClient:    userClient,
		jwtSecret:     jwtSecret,
		bcryptLimiter: bcryptLimiter,
	}
}

func (s *AuthServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	id, passwordHash, err := s.userClient.GetByEmail(ctx, email)
	if err != nil {
		if st, ok := status.FromError(err); ok && grpcutil.ShouldRetryCode(st.Code()) {
			return "", errors.New("user-service unavailable")
		}

		return "", errors.New("invalid credentials")
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
	// Ограничение на bcrypt
	s.bcryptLimiter <- struct{}{}
	defer func() { <-s.bcryptLimiter }()

	// Хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// Пытаемся создать пользователя
	id, err := s.userClient.Create(ctx, email, string(hash))
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.AlreadyExists {
				return 0, errors.New("email already registered")
			}
			if grpcutil.ShouldRetryCode(st.Code()) {
				return 0, errors.New("user-service unavailable")
			}
		}
		return 0, errors.New("internal error")
	}

	return id, nil
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

	// Явно проверяем срок жизни токена
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return 0, "", errors.New("token expired")
		}
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
