package ports

type Hub interface {
	Send(userID int64, event string, data any)
	Sockets() map[int64]struct{}
}
