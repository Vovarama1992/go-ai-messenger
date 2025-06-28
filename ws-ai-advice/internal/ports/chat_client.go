package ports

type ChatService interface {
	GetUserWithChatByThreadID(threadID string) (userID int64, chatID int64, email string, err error)
}
