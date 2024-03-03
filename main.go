package main

import (
	"gcom-backend/configs"
	"gcom-backend/controllers"
	"gcom-backend/docs"
	"gcom-backend/util"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			GCOM Backend
//	@version		1.0
//	@description	This is the backend service for UBC UAS

//	@contact.name	UBC UAS
//	@contact.url	https://ubcuas.com/
//	@contact.email	info@ubcuas.com

//	@host	localhost:1323

//	@Accept		json
//	@Produce	json
//	@Tags		Waypoints

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := configs.ConnectDatabase()

	mp, err := configs.ConnectMissionPlanner(os.Getenv("MP_URL"))
	if err != nil {
		log.Fatal("Error connecting to MPS")
	}

	e := echo.New()

	e.Use(util.DBMiddleware(db))
	e.Use(util.MPMiddleware(mp))
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, "Hello, World!")
	})

	//Swagger Docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//Waypoints
	e.POST("/waypoint", controllers.CreateWaypoint)
	e.POST("/waypoints", controllers.CreateWaypointBatch)
	e.PATCH("/waypoint/:waypointId", controllers.EditWaypoint)
	e.GET("/waypoint/:waypointId", controllers.GetWaypoint)
	e.DELETE("/waypoint/:waypointId", controllers.DeleteWaypoint)
	e.DELETE("/waypoints", controllers.DeleteWaypointBatch)
	e.GET("/waypoints", controllers.GetAllWaypoints)

	//Ground Object
	e.POST("/ground_object", controllers.CreateGroundObject)
	e.PATCH("/ground_object/:groundObjectId", controllers.EditGroundObject)
	e.GET("/ground_object/:groundObjectId", controllers.GetGroundObject)
	e.DELETE("/ground_object/:groundObjectId", controllers.DeleteGroundObject)

	//Payload
	e.POST("/payload", controllers.CreatePayload)
	e.PATCH("/payload:payloadId", controllers.EditPayload)
	e.GET("/payload:payloadId", controllers.GetPayload)
	e.DELETE("/payload/:payloadID", controllers.DeletePayload)

	//Drone
	e.GET("/status", controllers.GetCurrentStatus)
	e.GET("/status/history", controllers.GetStatusHistory)
	e.POST("/drone/takeoff", controllers.Takeoff)
	e.GET("/drone/land", controllers.Land)
	e.GET("/drone/rtl", controllers.RTL)
	e.GET("/drone/lock", controllers.Lock)
	e.GET("/drone/unlock", controllers.Unlock)
	e.GET("/drone/queue", controllers.GetQueue)
	e.POST("/drone/queue", controllers.PostQueue)
	e.POST("/drone/home", controllers.PostHome)

	//AirObjects
	e.GET("/air_object", controllers.GetAirObjects)
	e.DELETE("/air_object", controllers.DeleteAirObjects)
	e.POST("/air_object", controllers.CreateAirObjects)

	//Websockets
	e.Any("/socket.io/", controllers.WebsocketHandler())

	e.Logger.Fatal(e.Start("localhost:1323"))
}
