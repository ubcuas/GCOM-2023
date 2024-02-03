package controllers

import (
	"encoding/json"
	"gcom-backend/configs"
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
	return nil
}

func PostQueue(c echo.Context) error {
	return nil
}

