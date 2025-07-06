package model

import "time"

type Message struct {
	ID          int64     `json:"id"`
	ChatID      int64     `json:"chatId"`
	SenderID    int64     `json:"senderId"`
	Content     string    `json:"text"`
	AIGenerated bool      `json:"aiGenerated"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ThreadContext struct {
	UserID    int64
	ChatID    int64
	UserEmail string
	UserIDs   []int64
}
