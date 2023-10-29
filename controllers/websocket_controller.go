package controllers

import (
	"encoding/json"
	"gcom-backend/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{}

func DroneWebsocket(c echo.Context) error {
	var db = c.Get("db").(*gorm.DB)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	var drone models.Drone

	for {
		// Read received drone JSON
		_, msg, readErr := ws.ReadMessage()
		if readErr != nil {
			c.Logger().Error(err)
		}

		//Bind to struct
		if err := json.Unmarshal(msg, &drone); err != nil {
			err := ws.WriteMessage(websocket.TextMessage, []byte("FormatError"))
			if err != nil {
				c.Logger().Error(err)
			}
		} else {
			//Add to DB
			db.Save(&drone)

			//Delete drones older than 5 minutes
			var drones []models.Drone
			db.Delete(&drones, "timestamp < ?", time.Now().Unix()-300)
		}
	}
}
