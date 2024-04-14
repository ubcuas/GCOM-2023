package controllers

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zishang520/socket.io/v2/socket"
)

func WebsocketHandler() echo.HandlerFunc {
	io := socket.NewServer(nil, nil)

	io.On("connection", func(clients ...any) {
		fmt.Println("connected client")
		client := clients[0].(*socket.Socket)

		var telemetryTicker *time.Ticker
		var stopChannel chan bool

		client.On("ping", func(a ...any) {
			client.Emit("pong")
		})

		client.On("disconnect", func(a ...any) {
			fmt.Println("disconnected client")
		})

		client.On("telemetry_start", func(a ...any) {
			if telemetryTicker != nil {
				return
			}
			fmt.Println("telemetry_start")
			telemetryTicker = time.NewTicker(500 * time.Millisecond)
			stopChannel = make(chan bool)
			go func() {
				for {
					select {
					case <-telemetryTicker.C:
						client.Emit("telemetry_data", "DATAHERE")
					case <-stopChannel:
						return
					}
				}
			}()
		})

		client.On("telemetry_stop", func(a ...any) {
			if telemetryTicker == nil {
				return
			}
			fmt.Println("telemetry_stop")
			telemetryTicker.Stop()
			stopChannel <- true
			close(stopChannel)
			telemetryTicker = nil
		})
	})

	return func(c echo.Context) error {
		io.ServeHandler(nil).ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
