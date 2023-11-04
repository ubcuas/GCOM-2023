package util

import (
	"errors"
	engineio "github.com/googollee/go-engine.io"
	socketio "github.com/googollee/go-socket.io"
	"github.com/labstack/echo/v4"
)

/*
MIT License

Copyright (c) 2020 Maksim Pavlov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

type IWrapper interface {
	OnConnect(nsp string, f func(echo.Context, socketio.Conn) error)
	OnDisconnect(nsp string, f func(echo.Context, socketio.Conn, string))
	OnError(nsp string, f func(echo.Context, error))
	OnEvent(nsp, event string, f func(echo.Context, socketio.Conn, string))
	HandlerFunc(context echo.Context) error
}

type Wrapper struct {
	Context echo.Context
	Server  *socketio.Server
}

// NewWrapper Create wrapper and Socket.io server
func NewWrapper(options engineio.Options) (*Wrapper, error) {
	server := socketio.NewServer(options)

	return &Wrapper{
		Server: server,
	}, nil
}

// NewWrapperWithServer Create wrapper with exists Socket.io server
func NewWrapperWithServer(server *socketio.Server) (*Wrapper, error) {
	if server == nil {
		return nil, errors.New("socket.io server can not be nil")
	}

	return &Wrapper{
		Server: server,
	}, nil
}

// OnConnect On Socket.io client connect
func (s *Wrapper) OnConnect(nsp string, f func(echo.Context, socketio.Conn) error) {
	s.Server.OnConnect(nsp, func(conn socketio.Conn) error {
		return f(s.Context, conn)
	})
}

// OnDisconnect On Socket.io client disconnect
func (s *Wrapper) OnDisconnect(nsp string, f func(echo.Context, socketio.Conn, string)) {
	s.Server.OnDisconnect(nsp, func(conn socketio.Conn, msg string) {
		f(s.Context, conn, msg)
	})
}

// OnError On Socket.io error
func (s *Wrapper) OnError(nsp string, f func(echo.Context, error)) {
	s.Server.OnError(nsp, func(err error) {
		f(s.Context, err)
	})
}

// OnEvent On Socket.io event from client
func (s *Wrapper) OnEvent(nsp, event string, f func(echo.Context, socketio.Conn, string)) {
	s.Server.OnEvent(nsp, event, func(conn socketio.Conn, msg string) {
		f(s.Context, conn, msg)
	})
}

// HandlerFunc Handler function
func (s *Wrapper) HandlerFunc(context echo.Context) error {
	go s.Server.Serve()

	s.Context = context
	s.Server.ServeHTTP(context.Response(), context.Request())
	return nil
}
