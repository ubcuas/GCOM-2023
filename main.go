package main

import (
	"fmt"
	"gcom-backend/configs"
	"gcom-backend/controllers"
	_ "gcom-backend/docs"
	"gcom-backend/util"
	"log"
	"os"

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
	db := configs.Connect(false)
	err := os.MkdirAll("db", 0755)  //Create db dir
	err = os.MkdirAll("imgs", 0755) //Create images dir
	if err != nil {
		fmt.Println(err)
		return
	}

	mp, err := configs.ConnectMissionPlanner("http://host.docker.internal:9000")
	if err != nil {
		log.Fatal("Error connecting to MPS")
	}

	e := echo.New()
	e.Use(middleware.CORS())

	e.Use(util.DBMiddleware(db))
	e.Use(util.MPMiddleware(mp))
	e.Use(middleware.CORS())
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

	//Drone
	e.GET("/status", controllers.GetCurrentStatus)
	e.GET("/status/history", controllers.GetStatusHistory)
	e.POST("/drone/takeoff", controllers.Takeoff)
	e.GET("/drone/land", controllers.Land)
	e.POST("/drone/rtl", controllers.RTL)
	e.GET("/drone/lock", controllers.Lock)
	e.GET("/drone/unlock", controllers.Unlock)
	e.GET("/drone/queue", controllers.GetQueue)
	e.POST("/drone/queue", controllers.PostQueue)
	e.POST("/drone/home", controllers.PostHome)
	e.POST("/drone/arm", controllers.Arm)

	//Ground Objects
	e.POST("/groundobject", controllers.CreateGroundObject)
	e.POST("/groundobjects", controllers.CreateGroundObjectBatch)
	e.PATCH("/groundobject/:objectId", controllers.EditGroundObject)
	e.GET("/groundobject/:objectId", controllers.GetGroundObject)
	e.DELETE("/groundobject/:objectId", controllers.DeleteGroundObject)
	e.DELETE("/groundobjects", controllers.DeleteGroundObjectBatch)
	e.GET("/groundobjects", controllers.GetAllGroundObjects)

	//Image Handling
	e.POST("/image", controllers.UploadImage)
	e.GET("/image/list", controllers.ListImages)
	e.GET("/image/:filename", controllers.GetImage)

	//Websockets
	e.Any("/socket.io/", controllers.WebsocketHandler())

	e.Logger.Fatal(e.Start("0.0.0.0:1323"))
}
