package main

import (
	"gcom-backend/configs"
	"gcom-backend/controllers"
	_ "gcom-backend/docs"
	"gcom-backend/util"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"os"
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

	db := configs.Connect()
	mp, err := util.NewMissionPlanner(os.Getenv("mp_url"))
	if err != nil {
		log.Fatal("Error connecting to Mission Planner")
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
	e.PATCH("/waypoint/:waypointId", controllers.EditWaypoint)
	e.GET("/waypoint/:waypointId", controllers.GetWaypoint)
	e.DELETE("/waypoint/:waypointId", controllers.DeleteWaypoint)
	e.GET("/waypoints", controllers.GetAllWaypoints)

	//Drone
	e.GET("/status", controllers.GetCurrentStatus)
	e.GET("/status/history", controllers.GetStatusHistory)

	//Websockets
	e.Any("/socket.io/", controllers.WebsocketHandler())

	e.Logger.Fatal(e.Start(":1323"))
}
