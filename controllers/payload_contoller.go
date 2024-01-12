package controllers

import (
	"gcom-backend/models"
	"gcom-backend/responses"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CreatePayload creates a paylod
//
//	@Summary		Create a payload
//	@Description	Create a singular payload based on JSON, must have sentinel ID of "-1"
//	@Tags			Payload
//	@Accept			json
//	@Produce		json
//	@Param			object		body		models.Payload				true	"Payload Data"
//	@Success		200			{object}	responses.PayloadResponse	"Success"
//	@Failure		400			{object}	responses.ErrorResponse		"Invalid JSON or Payload Data"
//	@Failure		500			{object}	responses.ErrorResponse		"Internal Error Payload"
//	@Router			/payload [post]
func CreatePayload(c echo.Context) error {
	var payload models.Payload    //Declares an empty Payload class
	db, _ := c.Get("db").(*gorm.DB) //Obtains the DB instance from the context

	if err := c.Bind(&payload); err != nil {
		/*
			.Bind() basically tries to force the JSON data provided into the
			struct, using the json:"field" annotations we added to know what
			goes where. c.Bind directly edits the payload variable and as such
			we provide the memory address of it using the & symbol in front.
			c.Bind returns an error if something goes wrong, and we catch it here
			by checking if it is nil
		*/
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	if validationErr := validate.Struct(&payload); validationErr != nil {
		/*
			validate.Struct() validates the struct based on the validate
			annotation we provided in the struct definition
		*/
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
		/*
			Here we use GROM functions. GROM has already created tables for the
			model definitions we provided and knows what type `payload` is
		*/
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "An error occurred creating the payload"})
	}

	return c.JSON(http.StatusOK, responses.PayloadResponse{
		Message:  "Payload Created!",
		Payload: payload})
}