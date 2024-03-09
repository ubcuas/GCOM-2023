package controllers

import (
	"encoding/json"
	"gcom-backend/configs"
	"gcom-backend/models"
	"gcom-backend/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetCurrentStatus gets the current status of the drone
//
//	@Summary		Get drone status
//	@Description	Get the current status of the drone
//	@Tags			Drone
//	@Produce		json
//	@Success		200	{object}	models.Drone	"Success"
//	@Router			/status [get]
func GetCurrentStatus(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)
	drone := mp.GetStatus()
	return c.JSON(http.StatusOK, drone)
}

// GetStatusHistory gets the status of the drone for the last 5 minutes
//
//	@Summary		Get drone status history
//	@Description	Get drone status for the last 5 minutes
//	@Tags			Drone
//	@Produce		json
//	@Success		200	{object}	[]models.Drone	"Success"
//	@Router			/status/history [get]
func GetStatusHistory(c echo.Context) error {
	return c.JSON(http.StatusOK, c.Get("drone"))
}

func Takeoff(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)

	var altitude float64
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		return err
	} else {
		altitude = json_map["altitude"].(float64)
	}

	mp.Takeoff(altitude)
	return c.HTML(http.StatusAccepted, "")
}

func Land(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)
	if mp.Land() {
		return c.HTML(http.StatusAccepted, "")
	} else {
		return c.HTML(http.StatusInternalServerError, "")
	}
}

func RTL(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)
	if mp.ReturnHome() {
		return c.HTML(http.StatusAccepted, "")
	} else {
		return c.HTML(http.StatusInternalServerError, "")
	}
}

func Lock(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)
	if mp.Lock() {
		return c.HTML(http.StatusAccepted, "")
	} else {
		return c.HTML(http.StatusInternalServerError, "")
	}
}

func Unlock(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)
	if mp.Unlock() {
		return c.HTML(http.StatusAccepted, "")
	} else {
		return c.HTML(http.StatusInternalServerError, "")
	}
}

func GetQueue(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)
	var queue = mp.GetQueue()
	return c.JSON(http.StatusOK, queue)
}

func PostQueue(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)
	var queue []models.Waypoint
	if err := c.Bind(&queue); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	for i := 0; i < len(queue); i++ {
		if validationErr := validate.Struct(&(queue[i])); validationErr != nil {
			return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
				Message: "Invalid waypoints data",
				Data:    validationErr.Error()})
		}
	}

	if mp.SetQueue(queue) {
		return c.HTML(http.StatusAccepted, "")
	} else {
		return c.HTML(http.StatusInternalServerError, "")
	}
}

func PostHome(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)
	var wp models.Waypoint
	if err := c.Bind(&wp); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	if validationErr := validate.Struct(&wp); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid waypoint data",
			Data:    validationErr.Error()})
	}

	if mp.SetHome(wp) {
		return c.HTML(http.StatusAccepted, "")
	} else {
		return c.HTML(http.StatusInternalServerError, "")
	}
}
