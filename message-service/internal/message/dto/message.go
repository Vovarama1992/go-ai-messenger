package dto

type IncomingMessage struct {
	ChatID      int64  `json:"chatId"`
	SenderID    int64  `json:"senderId"`
	Text        string `json:"text"`
	AIGenerated bool   `json:"aiGenerated"`
}
