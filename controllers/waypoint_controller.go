package controllers

import (
	"errors"
	"gcom-backend/models"
	"gcom-backend/responses"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var validate = validator.New()

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
			Status:  http.StatusBadRequest,
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	if validationErr := validate.Struct(&waypoint); validationErr != nil {
		/*
			validate.Struct() validates the struct based on the validate
			annotation we provided in the struct definition
		*/
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid waypoint data",
			Data:    validationErr.Error()})
	}

	if createErr := db.Create(&waypoint).Error; createErr != nil {
		/*
			Here we use GROM functions. GROM has already created tables for the
			model definitions we provided and knows what type `waypoint` is
		*/
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "An error occurred creating the waypoint"})
	}

	return c.JSON(http.StatusOK, responses.WaypointResponse{
		Status:   http.StatusOK,
		Message:  "Waypoint Created!",
		Waypoint: waypoint,
	})
}

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
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    castErr.Error()})
	}

	if bindErr != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid JSON format",
			Data:    bindErr.Error()})
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
			Status:  http.StatusInternalServerError,
			Message: "An error occurred updating the waypoint",
			Data:    updateAction.Error.Error()})
	} else if updateAction.RowsAffected < 1 {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Status:  http.StatusNotFound,
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

	return c.JSON(http.StatusOK, responses.WaypointResponse{
		Status:   http.StatusOK,
		Message:  "Waypoint Updated!",
		Waypoint: updatedWaypoint,
	})
}

func GetWaypoint(c echo.Context) error {
	waypointId := c.Param("waypointId")
	var waypoint models.Waypoint
	db, _ := c.Get("db").(*gorm.DB)

	if err := db.First(&waypoint, waypointId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "No such waypoint exists!"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error whilst querying waypoint!"})
	}

	return c.JSON(http.StatusOK, responses.WaypointResponse{
		Status:   http.StatusOK,
		Message:  "Waypoint Found!",
		Waypoint: waypoint,
	})
}

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
			Status:  http.StatusInternalServerError,
			Message: "Error whilst deleting waypoint!"})
	} else if dbAction.RowsAffected < 1 {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "No such waypoint exists!"})
	}

	return c.JSON(http.StatusOK, responses.WaypointResponse{
		Status:   http.StatusOK,
		Message:  "Waypoint Found!",
		Waypoint: models.Waypoint{},
	})
}

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
			Status:  http.StatusInternalServerError,
			Message: "Error whilst querying waypoints!",
			Data:    err.Error()})
	}

	return c.JSON(http.StatusOK, responses.WaypointsResponse{
		Status:    http.StatusOK,
		Message:   "Waypoint Found!",
		Waypoints: waypoints,
	})

}
