package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gcom-backend/models"
	"gcom-backend/responses"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var validate = validator.New()

// CreateWaypoint creates a waypoint
//
//	@Summary		Create a waypoint
//	@Description	Create a singular waypoint based on JSON, must have sentinel ID of "-1"
//	@Tags			Waypoint
//	@Accept			json
//	@Produce		json
//	@Param			waypoint	body		models.Waypoint								true	"Waypoint Data"
//	@Success		200			{object}	responses.SingleResponse[models.Waypoint]	"Success"
//	@Failure		400			{object}	responses.ErrorResponse						"Invalid JSON or Waypoint Data"
//	@Failure		500			{object}	responses.ErrorResponse						"Internal Error Creating Waypoint"
//	@Router			/waypoint [post]
func CreateWaypoint(c echo.Context) error {
	var waypoint models.Waypoint    //Declares an empty Waypoint class
	db, _ := c.Get("db").(*gorm.DB) //Obtains the DB instance from the context

	if err := c.Bind(&waypoint); err != nil {
		/*
			.Bind() basically tries to force the JSON data provided into the
			struct, using the json:"field" annotations we added to know what
			goes where. c.Bind directly edits the waypoint variable and as such
			we provide the memory address of it using the & symbol in front.
			c.Bind returns an error if something goes wrong, and we catch it here
			by checking if it is nil
		*/
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	if validationErr := validate.Struct(&waypoint); validationErr != nil {
		/*
			validate.Struct() validates the struct based on the validate
			annotation we provided in the struct definition
		*/
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid waypoint data",
			Data:    validationErr.Error()})
	}

	if waypoint.ID != -1 {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Non-sentinel ID passed"})
	} else {
		waypoint.ID = 0
	}

	if createErr := db.Create(&waypoint).Error; createErr != nil {
		/*
			Here we use GROM functions. GROM has already created tables for the
			model definitions we provided and knows what type `waypoint` is
		*/
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "An error occurred creating the waypoint"})
	}

	return c.JSON(http.StatusOK, responses.SingleResponse[models.Waypoint]{
		Message: "Waypoint created!",
		Model:   waypoint})
}

// CreateWaypointBatch creates multiple waypoints
//
//	@Summary		Create multiple waypoints
//	@Description	Create multiple waypoints based on JSON, all must have sentinel ID of "-1"
//	@Tags			Waypoint
//	@Accept			json
//	@Produce		json
//	@Param			waypoints	body		[]models.Waypoint							true	"Array of Waypoint Data"
//	@Success		200			{object}	responses.MultipleResponse[models.Waypoint]	"Success"
//	@Failure		400			{object}	responses.ErrorResponse						"Invalid JSON or Waypoint Data"
//	@Failure		500			{object}	responses.ErrorResponse						"Internal Error Creating Waypoint"
//	@Router			/waypoints [post]
func CreateWaypointBatch(c echo.Context) error {
	var waypoints []models.Waypoint
	db, _ := c.Get("db").(*gorm.DB)

	if err := c.Bind(&waypoints); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	for i := 0; i < len(waypoints); i++ {
		if validationErr := validate.Struct(&(waypoints[i])); validationErr != nil {
			return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
				Message: "Invalid waypoints data",
				Data:    validationErr.Error()})
		}
		if waypoints[i].ID != -1 {
			return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
				Message: fmt.Sprintf("Non-sentinel ID passed for waypoint %d", i)})
		} else {
			waypoints[i].ID = 0
		}

	}

	if createErr := db.Create(&waypoints).Error; createErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "An error occurred creating the waypoints"})
	}

	return c.JSON(http.StatusOK, responses.MultipleResponse[models.Waypoint]{
		Message: "Waypoints created!",
		Models:  waypoints})
}

// EditWaypoint edits a waypoint
//
//	@Summary		Edit a waypoint
//	@Description	Edit a singular waypoint based on path param and JSON
//	@Tags			Waypoint
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int											true	"Waypoint ID"
//	@Param			fields	body		string										true	"JSON fields"	example({"name": "Whiskey})
//	@Success		200		{object}	responses.SingleResponse[models.Waypoint]	"Success"
//	@Failure		400		{object}	responses.ErrorResponse						"Invalid JSON or Waypoint ID"
//	@Failure		404		{object}	responses.ErrorResponse						"Waypoint Not Found"
//	@Failure		500		{object}	responses.ErrorResponse						"Internal Error Editing Waypoint"
//	@Router			/waypoint/{id} [patch]
func EditWaypoint(c echo.Context) error {
	/*
		This line obtains the parameter we specified in main.go, when we added
		`:waypointId` to the url.
	*/
	waypointStringId := c.Param("waypointId")
	var waypoint models.Waypoint
	db, _ := c.Get("db").(*gorm.DB)

	waypointId, castErr := strconv.Atoi(waypointStringId)
	bindErr := c.Bind(&waypoint)

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

	if waypoint.ID != 0 {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "ID is not editable"})
	}

	/*
		Here is a generic way to query for results, .Model() specifies the model
		(table) we wish to query from, and needs a pointer to some struct to
		know which model table it should get. .Where() takes the parts which
		follow after an SQL WHERE statement and uses ? to specify blanks, with
		the function filling in those ?'s with the arguments that follow in
		left-to-right order.
	*/
	updateAction := db.Model(&models.Waypoint{}).
		Where("id = ?", waypointId).
		Updates(&waypoint)

	if updateAction.Error != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "An error occurred updating the waypoint",
			Data:    updateAction.Error.Error()})
	} else if updateAction.RowsAffected < 1 {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Message: "No such waypoint exists!"})
	}

	/*
		Here, we are getting the updated waypoint information using a shorthand
		notation, which overwrites the waypoint variable with the information
		from the query, which is querying by the primaryKey specified in the
		second argument.
	*/
	var updatedWaypoint models.Waypoint
	db.First(&updatedWaypoint, waypointId)

	return c.JSON(http.StatusOK, responses.SingleResponse[models.Waypoint]{
		Message: "Waypoint updated!",
		Model:   updatedWaypoint,
	})
}

// GetWaypoint gets a waypoint
//
//	@Summary		Get a waypoint
//	@Description	Get a singular waypoint based on path param
//	@Tags			Waypoint
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int											true	"Waypoint ID"
//	@Success		200	{object}	responses.SingleResponse[models.Waypoint]	"Success"
//	@Failure		404	{object}	responses.ErrorResponse						"Waypoint Not Found"
//	@Failure		500	{object}	responses.ErrorResponse						"Internal Error Querying Waypoint"
//	@Router			/waypoint/{id} [get]
func GetWaypoint(c echo.Context) error {
	waypointId := c.Param("waypointId")
	var waypoint models.Waypoint
	db, _ := c.Get("db").(*gorm.DB)

	if err := db.First(&waypoint, waypointId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Message: "No such waypoint exists!"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Error whilst querying waypoint!"})
	}

	return c.JSON(http.StatusOK, responses.SingleResponse[models.Waypoint]{
		Message: "Waypoint found!",
		Model:   waypoint,
	})
}

// DeleteWaypoint deletes a waypoint
//
//	@Summary		Delete a waypoint
//	@Description	Delete a singular waypoint based on path param
//	@Tags			Waypoint
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int											true	"Waypoint ID"
//	@Success		200	{object}	responses.SingleResponse[models.Waypoint]	"Success (returns a blank Waypoint)"
//	@Failure		404	{object}	responses.ErrorResponse						"Waypoint Not Found"
//	@Failure		500	{object}	responses.ErrorResponse						"Internal Error Deleting Waypoint"
//	@Router			/waypoint/{id} [delete]
func DeleteWaypoint(c echo.Context) error {
	db, _ := c.Get("db").(*gorm.DB)
	waypointId := c.Param("waypointId")

	/*
		Because db.Delete() does not return a ErrNoRecordsFound error if we pass
		in an id that does not exist, we need to be a bit more creative to detect
		this. Here, we are checking if our action resulting in any rows changing
		and if not, telling the user that the waypoint did not exist anyway.
	*/
	dbAction := db.Delete(&models.Waypoint{}, waypointId)
	if err := dbAction.Error; err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Error whilst deleting waypoint!"})
	} else if dbAction.RowsAffected < 1 {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Message: "No requested waypoint exists!"})
	}

	return c.JSON(http.StatusOK, responses.SingleResponse[models.Waypoint]{
		Message: "Waypoint deleted!",
		Model:   models.Waypoint{},
	})
}

// DeleteWaypointBatch deletes multiple waypoints
//
//	@Summary		Delete multiple waypoints
//	@Description	Delete multiple waypoints based on json body
//	@Tags			Waypoint
//	@Accept			json
//	@Produce		json
//	@Param			ids	body		[]int										true	"Waypoint IDs"
//	@Success		200	{object}	responses.SingleResponse[models.Waypoint]	"Success (returns a blank Waypoint)"
//	@Failure		400	{object}	responses.ErrorResponse						"Invalid JSON or Waypoint IDs"
//	@Failure		404	{object}	responses.ErrorResponse						"Waypoints Not Found"
//	@Failure		500	{object}	responses.ErrorResponse						"Internal Error Deleting Waypoint"
//	@Router			/waypoints [delete]
func DeleteWaypointBatch(c echo.Context) error {
	db, _ := c.Get("db").(*gorm.DB)
	/*
		Instead of specifying one waypoint to delete through the uri param,
		multiple waypoint ids will be represented through json as an array of ints.
	*/
	body, _ := io.ReadAll(c.Request().Body) // not sure if reading the body into bytes needs err handling.

	// Direct unmarshalling of data; No binding overhead
	var waypointIDs []int
	if marshalErr := json.Unmarshal(body, &waypointIDs); marshalErr != nil {
		// Will error out due to invalid format, ex. [1,2,"a"]
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON/ID format",
			Data:    marshalErr.Error()})
	}

	// ID verification
	for _, id := range waypointIDs {
		// Negative id check
		if id < 0 {
			return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
				Message: "Invalid ID; Negative ID entered"})
		}
		// Check id exists before any deletion; prevents partial deletion
		var waypointTBValidated = models.Waypoint{}
		if err := db.First(&waypointTBValidated, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, responses.ErrorResponse{
				Message: fmt.Sprintf("Requested waypoint %d does not exist!", id),
				Data:    err.Error()})
		}
	}

	for _, id := range waypointIDs {
		dbAction := db.Delete(&models.Waypoint{}, id)
		// Theoretically shouldn't error out after validation, extra validation
		// must be implemented if db.Delete errors out.
		if err := dbAction.Error; err != nil {
			return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Message: fmt.Sprintf("Error whilst deleting waypoint with id %d", id)})
		}
	}

	return c.JSON(http.StatusOK, responses.SingleResponse[models.Waypoint]{
		Message: "Waypoints deleted!",
		Model:   models.Waypoint{},
	})
}

// GetAllWaypoints gets all waypoints in the database
//
//	@Summary		Get all waypoints
//	@Description	Get all waypoints in the database
//	@Tags			Waypoint
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.MultipleResponse[models.Waypoint]	"Success"
//	@Failure		500	{object}	responses.ErrorResponse						"Internal Error Querying Waypoints"
//	@Router			/waypoints [get]
func GetAllWaypoints(c echo.Context) error {
	var waypoints []models.Waypoint
	db, _ := c.Get("db").(*gorm.DB)

	/*
		If we use .Find() on an model array, it results all the models in the
		DB. echo.Map{} does the work of parsing the model array in to JSON for
		us here, as it has done with a single model before.
	*/
	if err := db.Find(&waypoints).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Error whilst querying waypoints!",
			Data:    err.Error()})
	}

	return c.JSON(http.StatusOK, responses.MultipleResponse[models.Waypoint]{
		Message: "Waypoints found!",
		Models:  waypoints,
	})
}
