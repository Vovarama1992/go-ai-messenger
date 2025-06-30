package dto

type IncomingMessage struct {
	ChatID      int64  `json:"chatId,omitempty"`
	SenderID    int64  `json:"senderId,omitempty"`
	Text        string `json:"text"`
	AIGenerated bool   `json:"aiGenerated"`
	ThreadID    string `json:"threadId,omitempty"`
}
