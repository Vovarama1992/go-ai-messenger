package auth

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	mocks "github.com/Vovarama1992/go-ai-messenger/auth-service/internal/mocks"
)

func getTestSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "test_secret"
	}
	return secret
}

func getTestBcryptLimiter() chan struct{} {
	return make(chan struct{}, 4) // В тестах пусть будет 4
}

// === СЦЕНАРИЙ 1: юзер не найден ===
func TestLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mocks.NewMockUserClient(ctrl)
	mock.EXPECT().
		GetByEmail(gomock.Any(), "test@example.com").
		Return(int64(0), "", errors.New("user not found"))

	service := NewAuthService(mock, getTestSecret(), getTestBcryptLimiter())
	_, err := service.Login(context.Background(), "test@example.com", "password123")

	if err == nil {
		t.Error("Ожидалась ошибка, но err == nil")
	}
}

// === СЦЕНАРИЙ 2: юзер найден, но пароль не подходит ===
func TestLogin_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash, _ := bcrypt.GenerateFromPassword([]byte("realpassword"), bcrypt.DefaultCost)

	mock := mocks.NewMockUserClient(ctrl)
	mock.EXPECT().
		GetByEmail(gomock.Any(), "test@example.com").
		Return(int64(1), string(hash), nil)

	service := NewAuthService(mock, getTestSecret(), getTestBcryptLimiter())
	_, err := service.Login(context.Background(), "test@example.com", "wrongpassword")

	if err == nil {
		t.Error("Ожидалась ошибка при неправильном пароле, но err == nil")
	}
}

// === СЦЕНАРИЙ 3: всё ок — получаем токен ===
func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash, _ := bcrypt.GenerateFromPassword([]byte("mypassword"), bcrypt.DefaultCost)

	mock := mocks.NewMockUserClient(ctrl)
	mock.EXPECT().
		GetByEmail(gomock.Any(), "test@example.com").
		Return(int64(1), string(hash), nil)

	service := NewAuthService(mock, getTestSecret(), getTestBcryptLimiter())
	token, err := service.Login(context.Background(), "test@example.com", "mypassword")

	if err != nil {
		t.Errorf("Не ожидалось ошибки: %v", err)
	}
	if token == "" {
		t.Error("Ожидался токен, но он пустой")
	}
}

func TestValidateToken_Success(t *testing.T) {
	service := NewAuthService(nil, getTestSecret(), getTestBcryptLimiter())

	claims := jwt.MapClaims{
		"sub":   float64(1),
		"email": "test@example.com",
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(getTestSecret()))

	userID, email, err := service.ValidateToken(context.Background(), signed)
	if err != nil {
		t.Errorf("Не ожидалась ошибка: %v", err)
	}
	if userID == 0 {
		t.Error("Ожидался userID > 0")
	}
	if email != "test@example.com" {
		t.Errorf("Ожидался email 'test@example.com', а получил %s", email)
	}
}

func TestValidateToken_InvalidFormat(t *testing.T) {
	service := NewAuthService(nil, getTestSecret(), getTestBcryptLimiter())

	_, _, err := service.ValidateToken(context.Background(), "foobar.invalid.token")
	if err == nil {
		t.Error("Ожидалась ошибка при невалидном токене")
	}
}

func TestValidateToken_MissingSub(t *testing.T) {
	service := NewAuthService(nil, getTestSecret(), getTestBcryptLimiter())

	claims := jwt.MapClaims{
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
		"email": "test@example.com",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(getTestSecret()))

	_, _, err := service.ValidateToken(context.Background(), signed)
	if err == nil {
		t.Error("Ожидалась ошибка при отсутствии sub")
	}
}
