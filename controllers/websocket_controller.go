package controllers

import (
	"encoding/json"
	"fmt"
	"gcom-backend/models"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zishang520/socket.io/v2/socket"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func WebsocketHandler() func(context echo.Context) error {
	io := socket.NewServer(nil, nil)
	db, err := gorm.Open(sqlite.Open("./db/database.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	io.On("connection", func(clients ...any) {
		fmt.Println("[SOCKET] Client Connected")
		client := clients[0].(*socket.Socket)

		client.On("ping", func(a ...any) {
			client.Emit("pong")
		})

		client.On("disconnect", func(a ...any) {
			fmt.Println("[SOCKET] Client Disconnected")
		})

		client.On("drone_update", func(a ...any) {
			// Front-end socket data event emitter
			client.Broadcast().Emit("fe_response", a[0])

			// Save drone data to database
			droneMap := a[0].(map[string]interface{})
			jsonString, err := json.Marshal(droneMap)
			if err != nil {
				client.Emit("error", err.Error())
			}

			var drone models.Drone
			// Read received drone JSON
			if err := json.Unmarshal(jsonString, &drone); err != nil {
				client.Emit("error", err.Error())
			} else {
				// fmt.Println(drone)
				//Add drone
				db.Save(&drone)
				//Delete drones older than 5 minutes
				var drones []models.Drone
				db.Delete(&drones, "timestamp < ?", time.Now().Unix()-300)
			}
		})

		// Front end manual request for data
		client.On("fe_request", func(a ...any) {
			var drone models.Drone
			db.Order("timestamp desc").First(&drone)
			client.Emit("fe_response", drone)
		})
	})

	return func(c echo.Context) error {
		io.ServeHandler(nil).ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
