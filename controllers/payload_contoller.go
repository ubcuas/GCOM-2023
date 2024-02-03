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

// CreatePayload creates a payload
//
//	@Summary		Create a payload
//	@Description	Create a singular payload based on JSON, must have sentinel ID of "-1"
//	@Tags			Payload
//	@Accept			json
//	@Produce		json
//	@Param			object	body		models.Payload				true	"Payload Data"
//	@Success		200		{object}	responses.PayloadResponse	"Success"
//	@Failure		400		{object}	responses.ErrorResponse		"Invalid JSON or Payload Data"
//	@Failure		500		{object}	responses.ErrorResponse		"Internal Error Payload"
//	@Router			/payload [post]
func CreatePayload(c echo.Context) error {
	var payload models.Payload
	db, _ := c.Get("db").(*gorm.DB)

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	if validationErr := validate.Struct(&payload); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid payload data",
			Data:    validationErr.Error()})
	}

	if payload.ID != -1 {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Non-sentinel ID passed"})
	} else {
		payload.ID = 0
	}

	if createErr := db.Create(&payload).Error; createErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "An error occurred creating the payload"})
	}

	return c.JSON(http.StatusOK, responses.PayloadResponse{
		Message: "Payload Created!",
		Payload: payload})
}

// EditPayload edits a payload
//
//	@Summary		Edit a payload.
//	@Description	Edit a singular payload based on path param and JSON
//	@Tags			Payload
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Payload ID"
//	@Param			fields	body		string						true	"JSON fields"	example({"name": "Whiskey"})
//	@Success		200		{object}	responses.PayloadResponse	"Success"
//	@Failure		400		{object}	responses.ErrorResponse		"Invalid JSON or Payload ID"
//	@Failure		404		{object}	responses.ErrorResponse		"Payload Not Found"
//	@Failure		500		{object}	responses.ErrorResponse		"Internal Error Editing Payload"
//	@Router			/payload/{id} [patch]
func EditPayload(c echo.Context) error {
	payloadStringId := c.Param("payloadId")
	var payload models.Payload
	db, _ := c.Get("db").(*gorm.DB)

	payloadId, castErr := strconv.Atoi(payloadStringId)
	bindErr := c.Bind(&payload)

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

	if payload.ID != 0 {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "ID is not editable"})
	}

	updateAction := db.Model(&models.Payload{}).
		Where("id = ?", payloadId).
		Updates(&payload)

	if updateAction.Error != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "An error occurred updating the payload",
			Data:    updateAction.Error.Error()})
	} else if updateAction.RowsAffected < 1 {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Message: "No such payload exists!"})
	}

	var updatedPayload models.Payload
	db.First(&updatedPayload, payloadId)

	return c.JSON(http.StatusOK, responses.PayloadResponse{
		Message: "Payload Updated!",
		Payload: updatedPayload,
	})
}

// GetPayload gets a payload by ID
//
//	@Summary		Get a payload by ID
//	@Description	Get a singular payload based on the provided ID
//	@Tags			Payload
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int							true	"Payload ID"
//	@Success		200	{object}	responses.PayloadResponse	"Success"
//	@Failure		404	{object}	responses.ErrorResponse		"Payload Not Found"
//	@Failure		500	{object}	responses.ErrorResponse		"Internal Error Querying Payload"
//	@Router			/payload/{id} [get]
func GetPayload(c echo.Context) error {
	payloadId := c.Param("payloadId")
	var payload models.Payload
	db, _ := c.Get("db").(*gorm.DB)

	if err := db.First(&payload, payloadId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Message: "No such payload exists!"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Error whilst querying payload!"})
	}

	return c.JSON(http.StatusOK, responses.PayloadResponse{
		Message: "Payload Found!",
		Payload: payload,
	})
}

// DeletePayload deletes a payload
//
//	@Summary		Delete a payload
//	@Description	Delete a singular payload based on path param
//	@Tags			Payload
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int							true	"Payload ID"
//	@Success		200	{object}	responses.PayloadResponse	"Success (returns a blank Payload)"
//	@Failure		404	{object}	responses.ErrorResponse		"Payload Not Found"
//	@Failure		500	{object}	responses.ErrorResponse		"Internal Error Deleting Payload"
//	@Router			/ground_object/{id} [delete]
func DeletePayload(c echo.Context) error {
	db, _ := c.Get("db").(*gorm.DB)
	payloadId := c.Param("payloadId")

	dbAction := db.Delete(&models.Payload{}, payloadId)
	if err := dbAction.Error; err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Error whilst deleting payload!"})
	} else if dbAction.RowsAffected < 1 {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Message: "No such payload exists!"})
	}

	return c.JSON(http.StatusOK, responses.PayloadResponse{
		Message: "Payload Deleted!",
		Payload: models.Payload{},
	})
}
