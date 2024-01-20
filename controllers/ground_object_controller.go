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
//	@Tags			Ground Object
//	@Accept			json
//	@Produce		json
//	@Param			object		body		models.GroundObject				true	"Ground Object Data"
//	@Success		200			{object}	responses.GroundObjectResponse	"Success"
//	@Failure		400			{object}	responses.ErrorResponse			"Invalid JSON or Ground Object Data"
//	@Failure		500			{object}	responses.ErrorResponse			"Internal Error Creating Ground Object"
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
		Message:  "Ground Object Created!",
		GroundObject: ground_object})
}

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
		Message:  "Ground Object Updated!",
		GroundObject: updatedGroundObject,
	})
}


func GetGroundObject(c echo.Context) error {
	groundObjectId := c.Param("groundObjectId")
	var groundObject models.GroundObject
	db, _ := c.Get("db").(*gorm.DB)

	if err := db.First(&groundObject, groundObjectId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Message: "No such groundObject exists!"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Error whilst querying ground object!"})
	}

	return c.JSON(http.StatusOK, responses.GroundObjectResponse{
		Message:  "Ground object Found!",
		GroundObject: groundObject,
	})
}