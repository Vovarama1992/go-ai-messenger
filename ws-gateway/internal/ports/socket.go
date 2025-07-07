package ports

type Conn interface {
	ID() string
	Close() error
	Emit(event string, data any) error
	Context() any
	SetContext(ctx any)
	Namespace() string
	RemoteHeader() map[string][]string
}
