package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
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
	return c.JSON(http.StatusOK, c.Get("drone"))
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
