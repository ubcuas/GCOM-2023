package controllers

import (
	"encoding/json"
	"gcom-backend/models"
	socketio "github.com/googollee/go-socket.io"
	"github.com/labstack/echo/v4"
	esi "github.com/umirode/echo-socket.io"
	"gorm.io/gorm"
	"time"
)

func WebsocketHandler(c echo.Context) error {
	wrapper, err := esi.NewWrapper(nil)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	wrapper.OnEvent("/drone", "info", DroneHandler)

	return wrapper.HandlerFunc(c)
}

func DroneHandler(c echo.Context, conn socketio.Conn, msg string) {
	var db = c.Get("db").(*gorm.DB)
	var drone models.Drone
	// Read received drone JSON
	if err := json.Unmarshal([]byte(msg), &drone); err != nil {
		conn.Emit("error", err.Error())
	} else {
		//Add to DB
		db.Save(&drone)
		//Delete drones older than 5 minutes
		var drones []models.Drone
		db.Delete(&drones, "timestamp < ?", time.Now().Unix()-300)
	}
}
