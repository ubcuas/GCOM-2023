package controllers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/zishang520/socket.io/v2/socket"
)

func WebsocketHandler() echo.HandlerFunc {
	io := socket.NewServer(nil, nil)

	io.On("connection", func(clients ...any) {
		fmt.Println("connected client")
		client := clients[0].(*socket.Socket)

		client.On("ping", func(a ...any) {
			client.Emit("pong")
		})

		client.On("disconnect", func(a ...any) {
			fmt.Println("disconnected client")
		})
	})

	return func(c echo.Context) error {
		io.ServeHandler(nil).ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
