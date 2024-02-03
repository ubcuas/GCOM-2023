package controllers

import (
	"gcom-backend/models"
	"gcom-backend/responses"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetAirObjects gets all present air objects (drones)
//
//	@Summary		Get air objects
//	@Description	Get the status of all current air objects (drones)
//	@Tags			AirObject
//	@Produce		json
//	@Success		200	{object}	responses.AirObjectsResponse	"Success"
//	@Failure		404	{object}	responses.ErrorResponse			"Objects not found"
//	@Failure		500	{object}	responses.ErrorResponse			"Internal error querying AirObjects"
//	@Router			/air_object [get]
func GetAirObjects(c echo.Context) error {
	var airObjects []models.AirObject
	db, _ := c.Get("db").(*gorm.DB)

	if err := db.Find(&airObjects).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Internal error querying AirObjects",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, responses.AirObjectsResponse{
		Message:    "AirObjects found!",
		AirObjects: airObjects,
	})
}

// DeleteAirObjects deletes all present air objects (remoteID drones)
//
//	@Summary		Delete all AirObjects
//	@Description	Delete all AirObjects (remoteID drones)
//	@Tags			AirObject
//	@Produce		json
//	@Success		200	{object}	responses.AirObjectsResponse	"Success (returns a empty array)"
//	@Failure		500	{object}	responses.ErrorResponse			"Internal Error Deleting AirObjects"
//	@Router			/air_object [delete]
func DeleteAirObjects(c echo.Context) error {
	db, _ := c.Get("db").(*gorm.DB)
	// gorm doesn't support global batch deleting without a where clause
	dbAction := db.Where("1 = 1").Delete(&models.AirObject{})
	if dbAction.Error != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Internal error deleting AirObjects",
			Data:    dbAction.Error.Error(),
		})
	}
	return c.JSON(http.StatusOK, responses.AirObjectsResponse{
		Message:    "All AirObjects deleted!",
		AirObjects: []models.AirObject{},
	})
}

// CreateAirObjects creates new AirObjects (remoteID drones) from a json array.
//
//	@Summary		Create multiple AirObjects
//	@Description	Create multiple AirObjects for remoteID based on JSON array.
//	@Tags			AirObject
//	@Accept			json
//	@Produce		json
//	@Param			airObjects	body		[]models.AirObject				true	"Array of AirObject Data"
//	@Success		200			{object}	responses.AirObjectsResponse	"Success"
//	@Failure		400			{object}	responses.ErrorResponse			"Invalid JSON or AirObject Data"
//	@Failure		500			{object}	responses.ErrorResponse			"Internal Error Creating AirObjects"
//	@Router			/air_object [post]
func CreateAirObjects(c echo.Context) error {
	var airObjects []models.AirObject
	db, _ := c.Get("db").(*gorm.DB)

	if err := c.Bind(&airObjects); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	for _, airObject := range airObjects {
		if validationErr := validate.Struct(airObject); validationErr != nil {
			return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
				Message: "Invalid AirObject format",
				Data:    validationErr.Error()})
		}
	}

	if createErr := db.Create(&airObjects).Error; createErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Internal error creating AirObjects",
			Data:    createErr.Error()})
	}

	return c.JSON(http.StatusOK, responses.AirObjectsResponse{
		Message:    "AirObjects created!",
		AirObjects: airObjects,
	})
}
