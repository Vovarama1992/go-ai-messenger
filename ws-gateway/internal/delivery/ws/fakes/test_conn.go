package ws

import (
	"net"
	"net/http"
	"net/url"
)

type FakeConn struct {
	Ctx any
}

func (f *FakeConn) ID() string                                      { return "" }
func (f *FakeConn) Namespace() string                               { return "" }
func (f *FakeConn) Emit(event string, args ...interface{})          {}
func (f *FakeConn) OnEvent(event string, handler interface{}) error { return nil }
func (f *FakeConn) Close() error                                    { return nil }
func (f *FakeConn) SetContext(ctx any)                              { f.Ctx = ctx }
func (f *FakeConn) Context() any                                    { return f.Ctx }
func (f *FakeConn) RemoteHeader() http.Header {
	return http.Header{}
}
func (f *FakeConn) URL() url.URL         { return url.URL{} }
func (f *FakeConn) LocalAddr() net.Addr  { return &net.IPAddr{} }
func (f *FakeConn) RemoteAddr() net.Addr { return &net.IPAddr{} }
func (f *FakeConn) Join(room string)     {}
func (f *FakeConn) Leave(room string)    {}
func (f *FakeConn) LeaveAll()            {}
func (f *FakeConn) Rooms() []string      { return nil }
