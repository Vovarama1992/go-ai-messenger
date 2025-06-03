package usecase_test

import (
	"context"
	"testing"

	"github.com/Vovarama1992/go-ai-messenger/user-service/internal/user/mocks"
	"github.com/Vovarama1992/go-ai-messenger/user-service/internal/user/model"
	"github.com/Vovarama1992/go-ai-messenger/user-service/internal/user/usecase"
	"go.uber.org/mock/gomock"
)

func TestGetByEmail_UserFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	expectedUser := &model.User{
		ID:    1,
		Email: "test@example.com",
	}

	mockRepo.EXPECT().
		FindByEmail(gomock.Any(), "test@example.com").
		Return(expectedUser, nil)

	s := usecase.NewUserService(mockRepo)

	user, err := s.GetByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Fatalf("unexpected result: %+v", user)
	}
}

func TestCreate_EmptyEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	s := usecase.NewUserService(mockRepo)

	_, err := s.Create(context.Background(), "", "some-hash")
	if err != usecase.ErrEmailRequired {
		t.Fatalf("ожидалась ошибка ErrEmailRequired, но получили: %v", err)
	}
}

func TestCreate_EmptyPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	s := usecase.NewUserService(mockRepo)

	_, err := s.Create(context.Background(), "test@example.com", "")
	if err != usecase.ErrPasswordRequired {
		t.Fatalf("ожидалась ошибка ErrPasswordRequired, но получили: %v", err)
	}
}

func TestCreate_EmailAlreadyTaken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().
		FindByEmail(gomock.Any(), "used@example.com").
		Return(&model.User{ID: 1, Email: "used@example.com"}, nil)

	s := usecase.NewUserService(mockRepo)

	_, err := s.Create(context.Background(), "used@example.com", "somehash")
	if err != usecase.ErrEmailTaken {
		t.Fatalf("ожидалась ошибка ErrEmailTaken, но получили: %v", err)
	}
}

func TestCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().
		FindByEmail(gomock.Any(), "new@example.com").
		Return(nil, nil)

	mockRepo.EXPECT().
		Create(gomock.Any(), "new@example.com", "somehash").
		Return(&model.User{ID: 2, Email: "new@example.com"}, nil)

	s := usecase.NewUserService(mockRepo)

	user, err := s.Create(context.Background(), "new@example.com", "somehash")
	if err != nil {
		t.Fatalf("не ожидалась ошибка: %v", err)
	}
	if user.Email != "new@example.com" {
		t.Fatalf("неверный результат: %+v", user)
	}
}
