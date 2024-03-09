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

// Takeoff tells the drone to take off to a specific altitude
//
//	@Summary		Take off Drone
//	@Description	Tells Drone to takeoff
//	@Tags			Drone
//	@Accept			json
//	@Param			altitude	body	number	true	"Takeoff Altitude"
//	@Success		200
//	@Router			/drone/takeoff [post]
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

// Land tells the drone to land
//
//	@Summary		Take off Drone
//	@Description	Tells Drone to land
//	@Tags			Drone
//	@Success		200	body	string	"Command issued successfully"
//	@Failure		500	body	string	"Command failed to be issued"
//	@Router			/drone/land [get]
func Land(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)
	if mp.Land() {
		return c.HTML(http.StatusAccepted, "")
	} else {
		return c.HTML(http.StatusInternalServerError, "")
	}
}

// RTL return to home waypoint and land
//
//	@Summary		Returns to Home and Lands
//	@Description	Tells Drone to return home and land
//	@Tags			Drone
//	@Success		200	body	string	"RTL command issued successfully"
//	@Failure		500	body	string	"RTL command encountered an error"
//	@Router			/drone/rtl [get]
func RTL(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)
	if mp.ReturnHome() {
		return c.HTML(http.StatusAccepted, "")
	} else {
		return c.HTML(http.StatusInternalServerError, "")
	}
}

// Lock locks the drone
//
//	@Summary		Halts drone in place while preserving queue
//	@Description	Stops drone movement while preserving existing queue
//	@Tags			Drone
//	@Success		200	body	string	"Drone locked successfully"
//	@Failure		500	body	string	"Drone unable to lock (already locked?)"
//	@Router			/drone/lock [get]
func Lock(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)
	if mp.Lock() {
		return c.HTML(http.StatusAccepted, "")
	} else {
		return c.HTML(http.StatusInternalServerError, "")
	}
}

// Unlock unlocks the drone
//
//	@Summary		Halts drone in place while preserving queue
//	@Description	Stops drone movement while preserving existing queue
//	@Tags			Drone
//	@Success		200	body	string	"Drone unlocked successfully"
//	@Failure		500	body	string	"Drone unable to unlock (already unlocked?)"
//	@Router			/drone/lock [get]
func Unlock(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)
	if mp.Unlock() {
		return c.HTML(http.StatusAccepted, "")
	} else {
		return c.HTML(http.StatusInternalServerError, "")
	}
}

// GetQueue obtains the current queue in MissionPlanner
//
//	@Summary		Returns queue in Mission Planner
//	@Description	Returns queue in Mission Planner
//	@Tags			Drone
//	@Produce		json
//	@Success		200	{object}	[]models.Waypoint
//	@Router			/drone/queue [get]
func GetQueue(c echo.Context) error {
	mp := c.Get("mp").(*configs.MissionPlanner)
	var queue = mp.GetQueue()
	return c.JSON(http.StatusOK, queue)
}

// PostQueue sends a queue to MissionPlanner
//
//	@Summary		Sends a queue in Mission Planner
//	@Description	Sends a queue in Mission Planner
//	@Tags			Drone
//	@Accept			json
//	@Param			waypoints	body	[]models.Waypoint	true	"Array of Waypoint Data"
//	@Success		200
//	@Router			/drone/queue [post]
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

// PostHome updates the home waypoint
//
//	@Summary		Updates the home waypoint
//	@Description	Updates the home waypoint
//	@Tags			Drone
//	@Accept			json
//	@Param			waypoints	body	models.Waypoint	true	"Home Waypoint"
//	@Success		200
//	@Router			/drone/home [post]
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
