package gpt

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/gptutil"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/ports"
	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	api         *openai.Client
	assistantID string
	apiKey      string // ← добавь это
	baseURL     string
}

var _ ports.GptClient = (*Client)(nil)

func NewClient(apiKey string, assistantID string, baseURL string) *Client {
	return &Client{
		api:         openai.NewClient(apiKey),
		assistantID: assistantID,
		apiKey:      apiKey,
		baseURL:     baseURL,
	}
}
func (c *Client) CreateThreadForUserAndChat(
	ctx context.Context,
	userEmail string,
	messages []dto.ChatMessage,
) (string, error) {
	thread, err := c.api.CreateThread(ctx, openai.ThreadRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to create thread: %w", err)
	}

	systemPrompt := fmt.Sprintf(
		"Ты создаёшь GPT thread для чата. Главный участник привязки: %s. "+
			"Вот история переписки между пользователями. Проанализируй её.",
		userEmail,
	)

	_, err = c.api.CreateMessage(ctx, thread.ID, openai.MessageRequest{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemPrompt,
	})
	if err != nil {
		return "", fmt.Errorf("failed to add system prompt: %w", err)
	}

	var sb strings.Builder
	for _, msg := range messages {
		sb.WriteString(fmt.Sprintf("%s: %s\n", msg.SenderEmail, msg.Text))
	}

	_, err = c.api.CreateMessage(ctx, thread.ID, openai.MessageRequest{
		Role:    openai.ChatMessageRoleUser,
		Content: sb.String(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to add full dialogue: %w", err)
	}

	return thread.ID, nil
}

func (c *Client) SendMessageToThread(
	ctx context.Context,
	threadID, role, content string,
) error {
	_, err := c.api.CreateMessage(ctx, threadID, openai.MessageRequest{
		Role:    role,
		Content: content,
	})
	if err != nil {
		return fmt.Errorf("failed to send message to thread %s: %w", threadID, err)
	}
	return nil
}

func (c *Client) SendMessageAndGetAutoreply(
	ctx context.Context,
	threadID string,
	senderEmail string,
	text string,
) (string, error) {
	prompt := fmt.Sprintf(
		"бро придумай ответ за юзера который привязал чат. сообщение пришло от: %s. текст сообщения: %s",
		senderEmail, text,
	)

	if _, err := c.api.CreateMessage(ctx, threadID, openai.MessageRequest{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	}); err != nil {
		return "", fmt.Errorf("failed to send message to thread: %w", err)
	}

	run, err := c.api.CreateRun(ctx, threadID, openai.RunRequest{
		AssistantID: c.assistantID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create run: %w", err)
	}

waitRun:
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			run, err = c.api.RetrieveRun(ctx, threadID, run.ID)
			if err != nil {
				return "", fmt.Errorf("failed to retrieve run: %w", err)
			}
			if run.Status == openai.RunStatusCompleted {
				break waitRun
			}
			if run.Status == openai.RunStatusFailed {
				return "", fmt.Errorf("run failed")
			}
			time.Sleep(1 * time.Second)
		}
	}

	return gptutil.GetLastAssistantMessage(ctx, c.apiKey, c.baseURL, threadID)
}

func (c *Client) GetAdvice(
	ctx context.Context,
	threadID string,
) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	if _, err := c.api.CreateMessage(ctx, threadID, openai.MessageRequest{
		Role:    openai.ChatMessageRoleUser,
		Content: "бро, дай совет пользователю по текущему состоянию чата",
	}); err != nil {
		return "", fmt.Errorf("failed to send advice request: %w", err)
	}

	run, err := c.api.CreateRun(ctx, threadID, openai.RunRequest{
		AssistantID: c.assistantID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create advice run: %w", err)
	}

waitRun:
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			run, err = c.api.RetrieveRun(ctx, threadID, run.ID)
			if err != nil {
				return "", fmt.Errorf("failed to retrieve run: %w", err)
			}
			if run.Status == openai.RunStatusCompleted {
				break waitRun
			}
			if run.Status == openai.RunStatusFailed {
				return "", fmt.Errorf("advice run failed")
			}
			time.Sleep(1 * time.Second)
		}
	}

	return gptutil.GetLastAssistantMessage(ctx, c.apiKey, c.baseURL, threadID)
}
