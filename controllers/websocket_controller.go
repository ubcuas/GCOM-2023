package controllers

import (
	"encoding/json"
	"fmt"
	"gcom-backend/models"
	"github.com/labstack/echo/v4"
	"github.com/zishang520/socket.io/v2/socket"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
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
				fmt.Println(drone)
				//Add drone
				db.Save(&drone)
				//Delete drones older than 5 minutes
				var drones []models.Drone
				db.Delete(&drones, "timestamp < ?", time.Now().Unix()-300)
			}
		})
	})

	return func(c echo.Context) error {
		io.ServeHandler(nil).ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
