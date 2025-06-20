package model

type ThreadResult struct {
	ChatID   int64  `json:"chatId"`
	UserID   int64  `json:"userId"`
	ThreadID string `json:"threadId"`
}
