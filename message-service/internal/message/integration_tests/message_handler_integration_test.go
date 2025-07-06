package integration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/delivery/grpc"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/infra/postgres"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/model"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/usecase"
	"github.com/Vovarama1992/go-ai-messenger/proto/messagepb"
)

type dummyUserClient struct{}

func (d *dummyUserClient) GetUserEmailByID(ctx context.Context, id int64) (string, error) {
	return "test@example.com", nil
}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:14",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "go_messenger",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithStartupTimeout(20 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	t.Cleanup(func() {
		_ = container.Terminate(ctx)
	})

	host, err := container.Host(ctx)
	assert.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432/tcp")
	assert.NoError(t, err)

	dsn := fmt.Sprintf("postgres://postgres:postgres@%s:%s/go_messenger?sslmode=disable", host, port.Port())
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		chat_id BIGINT,
		sender_id BIGINT,
		text TEXT,
		ai_generated BOOLEAN,
		created_at TIMESTAMP
	);`)
	assert.NoError(t, err)

	return db
}

func insertTestMessage(t *testing.T, repo *postgres.MessageRepo) int64 {
	msg := &model.Message{
		ChatID:      1001,
		SenderID:    5001,
		Content:     "integration hello",
		AIGenerated: true,
		CreatedAt:   time.Now(),
	}
	err := repo.Save(msg)
	assert.NoError(t, err)
	return msg.ID
}

func TestIntegration_GetMessagesByChat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewMessageRepo(db)
	user := &dummyUserClient{}
	service := usecase.NewMessageService(repo, user)
	handler := grpc.NewMessageHandler(service)

	insertTestMessage(t, repo)

	resp, err := handler.GetMessagesByChat(context.Background(), &messagepb.GetMessagesRequest{
		ChatId: 1001,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Messages)

	msg := resp.Messages[0]
	assert.Equal(t, int64(5001), msg.SenderId)
	assert.Equal(t, "integration hello", msg.Content)
	assert.Equal(t, "test@example.com", msg.SenderEmail)
}
