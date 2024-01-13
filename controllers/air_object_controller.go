package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetAirObjects gets all present air objects (drones)
//
//	@Summary		Get drone status
//	@Description	Get the current status of the drone
//	@Tags			Drone
//	@Produce		json
//	@Success		200	{object}	models.Drone	"Success"
//	@Router			/status [get]
func GetAirObjects(c echo.Context) error {
	return c.JSON(http.StatusOK, c.Get("drone"))
}
