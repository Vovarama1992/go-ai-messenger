package ws

import (
	"net"
	"net/http"
	"net/url"
)

type EmittedEvent struct {
	Event string
	Data  []interface{}
}

type FakeConnWithEmitLog struct {
	Ctx     any
	Emitted []EmittedEvent
}

func (f *FakeConnWithEmitLog) ID() string        { return "" }
func (f *FakeConnWithEmitLog) Namespace() string { return "" }
func (f *FakeConnWithEmitLog) Emit(event string, args ...interface{}) {
	f.Emitted = append(f.Emitted, EmittedEvent{Event: event, Data: args})
}
func (f *FakeConnWithEmitLog) OnEvent(event string, handler interface{}) error { return nil }
func (f *FakeConnWithEmitLog) Close() error                                    { return nil }
func (f *FakeConnWithEmitLog) SetContext(ctx any)                              { f.Ctx = ctx }
func (f *FakeConnWithEmitLog) Context() any                                    { return f.Ctx }
func (f *FakeConnWithEmitLog) RemoteHeader() http.Header                       { return http.Header{} }
func (f *FakeConnWithEmitLog) URL() url.URL                                    { return url.URL{} }
func (f *FakeConnWithEmitLog) LocalAddr() net.Addr                             { return &net.IPAddr{} }
func (f *FakeConnWithEmitLog) RemoteAddr() net.Addr                            { return &net.IPAddr{} }
func (f *FakeConnWithEmitLog) Join(room string)                                {}
func (f *FakeConnWithEmitLog) Leave(room string)                               {}
func (f *FakeConnWithEmitLog) LeaveAll()                                       {}
func (f *FakeConnWithEmitLog) Rooms() []string                                 { return nil }
