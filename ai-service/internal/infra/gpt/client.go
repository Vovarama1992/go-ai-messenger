package gpt

import (
	"context"
	"fmt"

	model "github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	api         *openai.Client
	assistantID string
}

func NewClient(apiKey string, assistantID string) *Client {
	return &Client{
		api:         openai.NewClient(apiKey),
		assistantID: assistantID,
	}
}

func (c *Client) CreateThreadForUserAndChat(
	ctx context.Context,
	userEmail string,
	messages []model.ChatMessage,
) (string, error) {
	// 1. Создаём пустой thread
	thread, err := c.api.CreateThread(ctx, openai.ThreadRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to create thread: %w", err)
	}

	// 2. Добавляем system message про "главного" пользователя
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

	// 3. Составляем единый блок всех сообщений
	fullDialogue := ""
	for _, msg := range messages {
		fullDialogue += fmt.Sprintf("%s: %s\n", msg.SenderEmail, msg.Text)
	}

	_, err = c.api.CreateMessage(ctx, thread.ID, openai.MessageRequest{
		Role:    openai.ChatMessageRoleUser,
		Content: fullDialogue,
	})
	if err != nil {
		return "", fmt.Errorf("failed to add full dialogue: %w", err)
	}

	return thread.ID, nil
}

func (c *Client) SendMessageToThread(
	ctx context.Context,
	threadID string,
	role string,
	content string,
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

	_, err := c.api.CreateMessage(ctx, threadID, openai.MessageRequest{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	run, err := c.api.CreateRun(ctx, threadID, openai.RunRequest{
		AssistantID: c.assistantID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create run: %w", err)
	}

	for {
		run, err = c.api.RetrieveRun(ctx, threadID, run.ID)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve run: %w", err)
		}
		if run.Status == openai.RunStatusCompleted {
			break
		}
		if run.Status == openai.RunStatusFailed {
			return "", fmt.Errorf("run failed")
		}
	}

	list, err := c.api.ListMessages(ctx, threadID, nil)
	if err != nil || len(list.Messages) == 0 {
		return "", fmt.Errorf("failed to list messages: %w", err)
	}
	return list.Messages[0].Content[0].Text.Value, nil
}

func (c *Client) GetAdvice(
	ctx context.Context,
	threadID string,
) (string, error) {
	// 1. Добавляем сообщение "бро, дай совет"
	_, err := c.api.CreateMessage(ctx, threadID, openai.MessageRequest{
		Role:    openai.ChatMessageRoleUser,
		Content: "бро, дай совет пользователю по текущему состоянию чата",
	})
	if err != nil {
		return "", fmt.Errorf("failed to send advice request: %w", err)
	}

	// 2. Запускаем Assistant
	run, err := c.api.CreateRun(ctx, threadID, openai.RunRequest{
		AssistantID: c.assistantID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create advice run: %w", err)
	}

	// 3. Ждём ответа
	for {
		run, err = c.api.RetrieveRun(ctx, threadID, run.ID)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve run: %w", err)
		}
		if run.Status == openai.RunStatusCompleted {
			break
		}
		if run.Status == openai.RunStatusFailed {
			return "", fmt.Errorf("advice run failed")
		}
	}

	// 4. Читаем последний ответ
	list, err := c.api.ListMessages(ctx, threadID, nil)
	if err != nil || len(list.Messages) == 0 {
		return "", fmt.Errorf("failed to list messages: %w", err)
	}
	return list.Messages[0].Content[0].Text.Value, nil
}
