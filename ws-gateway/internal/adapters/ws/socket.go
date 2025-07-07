package wsadapter

import socketio "github.com/googollee/go-socket.io"

type SocketConnAdapter struct {
	conn socketio.Conn
}

func NewSocketConnAdapter(c socketio.Conn) *SocketConnAdapter {
	return &SocketConnAdapter{conn: c}
}

func (s *SocketConnAdapter) ID() string {
	return s.conn.ID()
}

func (s *SocketConnAdapter) Close() error {
	return s.conn.Close()
}

func (s *SocketConnAdapter) Emit(event string, data any) error {
	s.conn.Emit(event, data)
	return nil
}

func (s *SocketConnAdapter) Context() any {
	return s.conn.Context()
}

func (s *SocketConnAdapter) SetContext(ctx any) {
	s.conn.SetContext(ctx)
}

func (s *SocketConnAdapter) Namespace() string {
	return s.conn.Namespace()
}

func (s *SocketConnAdapter) RemoteHeader() map[string][]string {
	return s.conn.RemoteHeader()
}
