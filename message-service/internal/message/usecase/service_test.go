package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/mock"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/model"
)

func TestMessageService_SaveMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockMessageRepo(ctrl)
	svc := NewMessageService(mockRepo, nil) // <- передаём nil для userClient

	msg := &model.Message{
		ChatID:      10,
		SenderID:    42,
		Text:        "hello world",
		AIGenerated: true,
	}

	mockRepo.
		EXPECT().
		Save(gomock.AssignableToTypeOf(&model.Message{})).
		Times(1).
		DoAndReturn(func(m *model.Message) error {
			assert.Equal(t, msg.ChatID, m.ChatID)
			assert.Equal(t, msg.SenderID, m.SenderID)
			assert.Equal(t, msg.Text, m.Text)
			assert.Equal(t, msg.AIGenerated, m.AIGenerated)
			assert.WithinDuration(t, time.Now(), m.CreatedAt, time.Second)
			m.ID = 777
			return nil
		})

	err := svc.SaveMessage(context.Background(), msg)

	assert.NoError(t, err)
	assert.Equal(t, int64(777), msg.ID)
}

func TestMessageService_SaveMessage_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockMessageRepo(ctrl)
	svc := NewMessageService(mockRepo, nil)

	msg := &model.Message{
		ChatID:   1,
		SenderID: 2,
		Text:     "fail case",
	}

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(errors.New("db error"))

	err := svc.SaveMessage(context.Background(), msg)

	assert.Error(t, err)
	assert.EqualError(t, err, "db error")
}
