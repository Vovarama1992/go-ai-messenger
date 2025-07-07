package ports

type Hub interface {
	Register(userID int64, conn Conn)
	Unregister(userID int64)
	JoinRoom(userID int64, chatID int64)
	LeaveRoom(userID int64, chatID int64)
	SendToRoom(chatID int64, event string, data any)
	HasConnection(userID int64) bool

	GetConn(userID int64) Conn
}
