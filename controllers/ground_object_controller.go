package controllers

import (
	"gcom-backend/models"
	"gcom-backend/responses"
	"net/http"

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
	var ground_object models.GroundObject    //Declares an empty Ground Object class
	db, _ := c.Get("db").(*gorm.DB) //Obtains the DB instance from the context

	if err := c.Bind(&ground_object); err != nil {
		/*
			.Bind() basically tries to force the JSON data provided into the
			struct, using the json:"field" annotations we added to know what
			goes where. c.Bind directly edits the ground object variable and as such
			we provide the memory address of it using the & symbol in front.
			c.Bind returns an error if something goes wrong, and we catch it here
			by checking if it is nil
		*/
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	if validationErr := validate.Struct(&ground_object); validationErr != nil {
		/*
			validate.Struct() validates the struct based on the validate
			annotation we provided in the struct definition
		*/
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
		/*
			Here we use GROM functions. GROM has already created tables for the
			model definitions we provided and knows what type `ground_object` is
		*/
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "An error occurred creating the ground object"})
	}

	return c.JSON(http.StatusOK, responses.GroundObjectResponse{
		Message:  "Ground Object Created!",
		GroundObject: ground_object})
}