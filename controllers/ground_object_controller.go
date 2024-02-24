package controllers

import (
	"errors"
	"gcom-backend/models"
	"gcom-backend/responses"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CreateGroundObject creates a ground object
//
//	@Summary		Create a ground object
//	@Description	Create a singular ground object based on JSON, must have sentinel ID of "-1"
//	@Tags			GroundObject
//	@Accept			json
//	@Produce		json
//	@Param			object	body		models.GroundObject				true	"Ground Object Data"
//	@Success		200		{object}	responses.GroundObjectResponse	"Success"
//	@Failure		400		{object}	responses.ErrorResponse			"Invalid JSON or Ground Object Data"
//	@Failure		500		{object}	responses.ErrorResponse			"Internal Error Creating Ground Object"
//	@Router			/ground_object [post]
func CreateGroundObject(c echo.Context) error {
	var ground_object models.GroundObject
	db, _ := c.Get("db").(*gorm.DB)

	if err := c.Bind(&ground_object); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	if validationErr := validate.Struct(&ground_object); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid ground object data",
			Data:    validationErr.Error()})
	}

	if ground_object.ID != -1 {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Non-sentinel ID passed"})
	} else {
		ground_object.ID = 0
	}

	if createErr := db.Create(&ground_object).Error; createErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "An error occurred creating the ground object"})
	}

	return c.JSON(http.StatusOK, responses.GroundObjectResponse{
		Message:      "Ground Object Created!",
		GroundObject: ground_object})
}

// EditGroundObject edits a ground object
//
//	@Summary		Edit a ground object.
//	@Description	Edit a singular ground object based on path param and JSON
//	@Tags			GroundObject
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"GroundObject ID"
//	@Param			fields	body		string							true	"JSON fields"	example({"name": "Whiskey"})
//	@Success		200		{object}	responses.GroundObjectResponse	"Success"
//	@Failure		400		{object}	responses.ErrorResponse			"Invalid JSON or GroundObject ID"
//	@Failure		404		{object}	responses.ErrorResponse			"GroundObject Not Found"
//	@Failure		500		{object}	responses.ErrorResponse			"Internal Error Editing GroundObject"
//	@Router			/ground_object/{id} [patch]
func EditGroundObject(c echo.Context) error {
	groundObjectStringId := c.Param("groundObjectId")
	var groundObject models.GroundObject
	db, _ := c.Get("db").(*gorm.DB)

	groundObjectId, castErr := strconv.Atoi(groundObjectStringId)
	bindErr := c.Bind(&groundObject)

	if castErr != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid ID",
			Data:    castErr.Error()})
	}

	if bindErr != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    bindErr.Error()})
	}

	if groundObject.ID != 0 {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "ID is not editable"})
	}

	updateAction := db.Model(&models.GroundObject{}).
		Where("id = ?", groundObjectId).
		Updates(&groundObject)

	if updateAction.Error != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "An error occurred updating the ground object",
			Data:    updateAction.Error.Error()})
	} else if updateAction.RowsAffected < 1 {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Message: "No such ground object exists!"})
	}

	var updatedGroundObject models.GroundObject
	db.First(&updatedGroundObject, groundObjectId)

	return c.JSON(http.StatusOK, responses.GroundObjectResponse{
		Message:      "Ground Object Updated!",
		GroundObject: updatedGroundObject,
	})
}

// GetGroundObject gets a ground object by ID
//
//	@Summary		Get a ground object by ID
//	@Description	Get a singular ground object based on the provided ID
//	@Tags			GroundObject
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int								true	"Ground Object ID"
//	@Success		200	{object}	responses.GroundObjectResponse	"Success"
//	@Failure		404	{object}	responses.ErrorResponse			"Ground Object Not Found"
//	@Failure		500	{object}	responses.ErrorResponse			"Internal Error Querying Ground Object"
//	@Router			/ground_object/{id} [get]
func GetGroundObject(c echo.Context) error {
	groundObjectId := c.Param("groundObjectId")
	var groundObject models.GroundObject
	db, _ := c.Get("db").(*gorm.DB)

	if err := db.First(&groundObject, groundObjectId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Message: "No such ground object exists!"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Error whilst querying ground object!"})
	}

	return c.JSON(http.StatusOK, responses.GroundObjectResponse{
		Message:      "Ground object Found!",
		GroundObject: groundObject,
	})
}

// DeleteGroundObject deletes a ground object
//
//	@Summary		Delete a ground object
//	@Description	Delete a singular ground object based on path param
//	@Tags			GroundObject
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int								true	"GroundObject ID"
//	@Success		200	{object}	responses.GroundObjectResponse	"Success (returns a blank Ground)"
//	@Failure		404	{object}	responses.ErrorResponse			"GroundObject Not Found"
//	@Failure		500	{object}	responses.ErrorResponse			"Internal Error Deleting GroundObject"
//	@Router			/ground_object/{id} [delete]
func DeleteGroundObject(c echo.Context) error {
	db, _ := c.Get("db").(*gorm.DB)
	groundObjectId := c.Param("groundObjectId")

	dbAction := db.Delete(&models.GroundObject{}, groundObjectId)
	if err := dbAction.Error; err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Error whilst deleting ground object!"})
	} else if dbAction.RowsAffected < 1 {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Message: "No such ground object exists!"})
	}

	return c.JSON(http.StatusOK, responses.GroundObjectResponse{
		Message:      "Ground object Deleted!",
		GroundObject: models.GroundObject{},
	})
}
