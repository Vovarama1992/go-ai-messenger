package model

type EnrichedAdvice struct {
	UserID int64  `json:"-"`
	ChatID int64  `json:"chatId"`
	Text   string `json:"text"`
}
