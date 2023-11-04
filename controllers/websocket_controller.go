package controllers

import (
	"encoding/json"
	"fmt"
	"gcom-backend/models"
	esi "gcom-backend/util"
	socketio "github.com/googollee/go-socket.io"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"time"
)

func WebsocketHandler() func(context echo.Context) error {
	wrapper, err := esi.NewWrapper()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	wrapper.OnConnect("", func(context echo.Context, conn socketio.Conn) error {
		conn.SetContext("")
		context.Logger().Infof("SocketIO: Client Connected (ID: ", conn.ID(), ")")
		return nil
	})
	wrapper.OnError("", func(context echo.Context, e error) {
		context.Logger().Infof("SocketIO: ", e)
	})

	wrapper.OnDisconnect("", func(context echo.Context, conn socketio.Conn, msg string) {
		context.Logger().Infof("SocketIO: Client Disconnected ( ", msg, ")")
	})

	wrapper.OnEvent("/drone", "update", DroneHandler)

	return wrapper.HandlerFunc
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
