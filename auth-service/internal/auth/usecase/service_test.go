package auth

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

type MockUserClient struct {
	GetByEmailFunc func(ctx context.Context, email string) (int64, string, error)
	CreateFunc     func(ctx context.Context, email, passwordHash string) (int64, error)
}

func (m *MockUserClient) GetByEmail(ctx context.Context, email string) (int64, string, error) {
	return m.GetByEmailFunc(ctx, email)
}

func (m *MockUserClient) Create(ctx context.Context, email, passwordHash string) (int64, error) {
	return m.CreateFunc(ctx, email, passwordHash)
}

func getTestSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "test_secret"
	}
	return secret
}

// === СЦЕНАРИЙ 1: юзер не найден ===
func TestLogin_UserNotFound(t *testing.T) {
	mock := &MockUserClient{
		GetByEmailFunc: func(ctx context.Context, email string) (int64, string, error) {
			return 0, "", errors.New("user not found")
		},
	}

	service := NewAuthService(mock, getTestSecret())
	_, err := service.Login(context.Background(), "test@example.com", "password123")

	if err == nil {
		t.Error("Ожидалась ошибка, но err == nil")
	}
}

// === СЦЕНАРИЙ 2: юзер найден, но пароль не подходит ===
func TestLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("realpassword"), bcrypt.DefaultCost)

	mock := &MockUserClient{
		GetByEmailFunc: func(ctx context.Context, email string) (int64, string, error) {
			return 1, string(hash), nil
		},
	}

	service := NewAuthService(mock, getTestSecret())
	_, err := service.Login(context.Background(), "test@example.com", "wrongpassword")

	if err == nil {
		t.Error("Ожидалась ошибка при неправильном пароле, но err == nil")
	}
}

// === СЦЕНАРИЙ 3: всё ок — получаем токен ===
func TestLogin_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("mypassword"), bcrypt.DefaultCost)

	mock := &MockUserClient{
		GetByEmailFunc: func(ctx context.Context, email string) (int64, string, error) {
			return 1, string(hash), nil
		},
	}

	service := NewAuthService(mock, getTestSecret())
	token, err := service.Login(context.Background(), "test@example.com", "mypassword")

	if err != nil {
		t.Errorf("Не ожидалось ошибки: %v", err)
	}
	if token == "" {
		t.Error("Ожидался токен, но он пустой")
	}
}

func TestValidateToken_Success(t *testing.T) {
	service := NewAuthService(nil, getTestSecret())

	// сгенерим валидный токен
	token, err := service.Login(context.Background(), "test@example.com", "testpass")
	if err != nil {
		t.Fatalf("Не удалось создать токен: %v", err)
	}

	userID, err := service.ValidateToken(context.Background(), token)
	if err != nil {
		t.Errorf("Не ожидалась ошибка: %v", err)
	}
	if userID == 0 {
		t.Error("Ожидался userID > 0")
	}
}

func TestValidateToken_InvalidFormat(t *testing.T) {
	service := NewAuthService(nil, getTestSecret())

	_, err := service.ValidateToken(context.Background(), "foobar.invalid.token")
	if err == nil {
		t.Error("Ожидалась ошибка при невалидном токене")
	}
}

func TestValidateToken_MissingSub(t *testing.T) {
	service := NewAuthService(nil, getTestSecret())

	claims := jwt.MapClaims{
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(getTestSecret()))

	_, err := service.ValidateToken(context.Background(), signed)
	if err == nil {
		t.Error("Ожидалась ошибка при отсутствии sub")
	}
}
