package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gcom-backend/models"
	"gcom-backend/responses"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Timestamp int64   `json:"timestamp"`
}

// CreateGroundObject creates a ground object
//
//	@Summary		Create a ground object
//	@Description	Create a singular ground object based on JSON, must have sentinel ID of "-1"
//	@Tags			GroundObject
//	@Accept			json
//	@Produce		json
//	@Param			object	body		models.GroundObject								true	"Ground Object Data"
//	@Success		200		{object}	responses.SingleResponse[models.GroundObject]	"Success"
//	@Failure		400		{object}	responses.ErrorResponse							"Invalid JSON or Object Data"
//	@Failure		500		{object}	responses.ErrorResponse							"Internal Error Creating GroundObject"
//	@Router			/groundobject [post]
func CreateGroundObject(c echo.Context) error {
	var object models.GroundObject
	db, _ := c.Get("db").(*gorm.DB)

	if err := c.Bind(&object); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	if validationErr := validate.Struct(&object); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid object data",
			Data:    validationErr.Error()})
	}

	if object.ID != -1 {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Non-sentinel ID passed"})
	} else {
		object.ID = 0
	}

	if createErr := db.Create(&object).Error; createErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "An error occurred creating the object"})
	}

	return c.JSON(http.StatusOK, responses.SingleResponse[models.GroundObject]{
		Message: "GroundObject created!",
		Model:   object})
}

// getDubinsAndNotifyPayload creates new payload path and notifies drone
//
//	@Summary
//	@Description
//	@Tags
//	@Accept			json
//	@Produce		json
//	@Param			object	body		models.GroundObject								true	"Ground Object Data"
//	@Success		200		{object}	responses.SingleResponse[models.GroundObject]	"Success"
//	@Failure		400		{object}	responses.ErrorResponse							"Invalid JSON or Object Data"
//	@Failure		500		{object}	responses.ErrorResponse							"Internal Error Creating GroundObject"
//	@Router			/odlc-found [post]
func getDroneDataAtTimestamp(db *gorm.DB, timestamp int64) (*models.Drone, error) {
	var data models.Drone
	// Use Where() to specify the condition and First() to retrieve the first record that matches the condition
	result := db.Where("timestamp = ?", timestamp).First(&data)
	fmt.Printf("RESULT: %+v\n", data)
	if result.Error != nil {
		return nil, result.Error
	}

	return &data, nil
}

func GetDubinsAndNotifyPayload(c echo.Context) error {

	db, _ := c.Get("db").(*gorm.DB)

	// Create a new instance of Location
	loc := new(Location)
	initial_body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Error reading request body",
			Data:    err.Error()})
	}
	fmt.Printf("Request body: %s\n", initial_body)

	if err := json.Unmarshal(initial_body, loc); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	// Print the latitude, longitude, and altitude
	fmt.Printf("Latitude: %f, Longitude: %f, time: %d\n", loc.Latitude, loc.Longitude, loc.Timestamp)

	droneData, err := getDroneDataAtTimestamp(db, loc.Timestamp)
	if err != nil {
		log.Print(err)
	} else {
		fmt.Printf("RETRIEVED DATA: %+v\n", droneData)
		loc.Latitude = droneData.Latitude
		loc.Longitude = droneData.Longitude
	}
	fmt.Printf("Target Data: %+v\n", loc)

	// Create an instance of the Drone struct
	var drone models.Drone
	// use db to get drone data at timestamp

	// Unmarshal the JSON data into the Drone struct
	// err := json.Unmarshal(jsonDataByte, &drone)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
	// 		Message: "Invalid JSON format",
	// 		Data:    err.Error()})

	// }

	// db, _ := c.Get("db").(*gorm.DB)\
	db.Order("timestamp desc").First(&drone)
	// waypoints := []models.Waypoint{}
	/*
		If we use .Find() on an model array, it results all the models in the
		DB. echo.Map{} does the work of parsing the model array in to JSON for
		us here, as it has done with a single model before.
	*/
	// if err := db.Find(&waypoints).Error; err != nil {
	// 	return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
	// 		Message: "Error whilst querying waypoints!",
	// 		Data:    err.Error()})
	// }
	// nextWaypoint := models.Waypoint{
	// 	ID:        0,
	// 	Name:      "Test",
	// 	Latitude:  drone.Latitude,
	// 	Longitude: drone.Longitude,
	// 	Altitude:  drone.Altitude,
	// }
	// if len(waypoints) != 0 {
	// 	nextWaypoint = waypoints[0]
	// }

	// Sample JSON data (can be replaced with GetCurrentStatus() result)
	// jsonData := controllers.getCurrentStatus()

	payload := map[string]interface{}{
		"current_waypoint": map[string]interface{}{
			"id":        -1,
			"name":      "",
			"latitude":  drone.Latitude,
			"longitude": drone.Longitude,
			"altitude":  drone.Altitude,
		},
		"desired_waypoint": map[string]interface{}{
			"id":        -1,
			"name":      "",
			"latitude":  loc.Latitude,
			"longitude": loc.Longitude,
			"altitude":  25,
		},
		"current_heading": drone.Heading,
		"desired_heading": drone.Heading,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error:", err)
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	// Send POST request to localhost:7001
	postResp, err := http.Post("http://localhost:7010", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}
	defer postResp.Body.Close()

	// Read the response body
	postBody, err := io.ReadAll(postResp.Body)
	if err != nil {
		fmt.Println("Error reading POST response body:", err)
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}
	fmt.Println("POST Response:", string(postBody))

	// TODO: send payload to drone
	// UPDATE THIS
	url := "http://192.168.1.65:9000/insert"
	jsonBody := bytes.NewBuffer(postBody)
	resp2, err := http.Post(url, "application/json", jsonBody)

	if err != nil {
		fmt.Println("Error sending POST request:", err)
	}
	print("MPS Response: ", resp2)

	return c.JSON(http.StatusOK, "worked!!!!")
}

// CreateGroundObjectBatch creates multiple ground objects
//
//	@Summary		Create multiple ground objects
//	@Description	Create multiple ground objects based on JSON, all must have sentinel ID of "-1"
//	@Tags			GroundObject
//	@Accept			json
//	@Produce		json
//	@Param			objects	body		[]models.GroundObject							true	"Array of object Data"
//	@Success		200		{object}	responses.MultipleResponse[models.GroundObject]	"Success"
//	@Failure		400		{object}	responses.ErrorResponse							"Invalid JSON or object Data"
//	@Failure		500		{object}	responses.ErrorResponse							"Internal Error Creating GroundObject"
//	@Router			/groundobjects [post]
func CreateGroundObjectBatch(c echo.Context) error {
	var objects []models.GroundObject
	db, _ := c.Get("db").(*gorm.DB)

	if err := c.Bind(&objects); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON format",
			Data:    err.Error()})
	}

	for i := 0; i < len(objects); i++ {
		if validationErr := validate.Struct(&(objects[i])); validationErr != nil {
			return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
				Message: "Invalid objects data",
				Data:    validationErr.Error()})
		}
		if objects[i].ID != -1 {
			return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
				Message: fmt.Sprintf("Non-sentinel ID passed for object %d", i)})
		} else {
			objects[i].ID = 0
		}

	}

	if createErr := db.Create(&objects).Error; createErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "An error occurred creating the ground object"})
	}

	return c.JSON(http.StatusOK, responses.MultipleResponse[models.GroundObject]{
		Message: "GroundObjects created!",
		Models:  objects})
}

// EditGroundObject edits a ground object
//
//	@Summary		Edit a ground object
//	@Description	Edit a singular object based on path param and JSON
//	@Tags			GroundObject
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int												true	"GroundObject ID"
//	@Param			fields	body		string											true	"JSON fields"	example({"color": "black"})
//	@Success		200		{object}	responses.SingleResponse[models.GroundObject]	"Success"
//	@Failure		400		{object}	responses.ErrorResponse							"Invalid JSON or GroundObject ID"
//	@Failure		404		{object}	responses.ErrorResponse							"GroundObject Not Found"
//	@Failure		500		{object}	responses.ErrorResponse							"Internal Error Editing GroundObject"
//	@Router			/groundobject/{id} [patch]
func EditGroundObject(c echo.Context) error {
	objectStringId := c.Param("objectId")
	var object models.GroundObject
	db, _ := c.Get("db").(*gorm.DB)

	objectId, castErr := strconv.Atoi(objectStringId)
	bindErr := c.Bind(&object)

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

	if object.ID != 0 {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "ID is not editable"})
	}

	updateAction := db.Model(&models.GroundObject{}).
		Where("id = ?", objectId).
		Updates(&object)

	if updateAction.Error != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "An error occurred updating the object",
			Data:    updateAction.Error.Error()})
	} else if updateAction.RowsAffected < 1 {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Message: "No such object exists!"})
	}

	var updatedObject models.GroundObject
	db.First(&updatedObject, objectId)

	return c.JSON(http.StatusOK, responses.SingleResponse[models.GroundObject]{
		Message: "GroundObject updated!",
		Model:   updatedObject,
	})
}

// GetGroundObject gets a ground object
//
//	@Summary		Get a ground object
//	@Description	Get a singular ground object based on path param
//	@Tags			GroundObject
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int												true	"Ground Object ID"
//	@Success		200	{object}	responses.SingleResponse[models.GroundObject]	"Success"
//	@Failure		404	{object}	responses.ErrorResponse							"Object Not Found"
//	@Failure		500	{object}	responses.ErrorResponse							"Internal Error Querying GroundObject"
//	@Router			/groundobject/{id} [get]
func GetGroundObject(c echo.Context) error {
	objectId := c.Param("objectId")
	var groundObject models.GroundObject
	db, _ := c.Get("db").(*gorm.DB)

	if err := db.First(&groundObject, objectId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Message: "No such object exists!"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Error whilst querying object!"})
	}

	return c.JSON(http.StatusOK, responses.SingleResponse[models.GroundObject]{
		Message: "Object found!",
		Model:   groundObject,
	})
}

// DeleteGroundObject deletes a ground object
//
//	@Summary		Delete a ground object
//	@Description	Delete a singular ground object based on path param
//	@Tags			GroundObject
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int												true	"GroundObject ID"
//	@Success		200	{object}	responses.SingleResponse[models.GroundObject]	"Success (returns a blank GroundObject)"
//	@Failure		404	{object}	responses.ErrorResponse							"GroundObject Not Found"
//	@Failure		500	{object}	responses.ErrorResponse							"Internal Error Deleting GroundObject"
//	@Router			/groundobject/{id} [delete]
func DeleteGroundObject(c echo.Context) error {
	db, _ := c.Get("db").(*gorm.DB)
	objectId := c.Param("objectId")
	dbAction := db.Delete(&models.GroundObject{}, objectId)
	if err := dbAction.Error; err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Error whilst deleting object!"})
	} else if dbAction.RowsAffected < 1 {
		return c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Message: "No requested object exists!"})
	}

	return c.JSON(http.StatusOK, responses.SingleResponse[models.GroundObject]{
		Message: "GroundObject deleted!",
		Model:   models.GroundObject{},
	})
}

// DeleteGroundObjectBatch deletes multiple ground objects
//
//	@Summary		Delete multiple ground objects
//	@Description	Delete multiple ground objects based on json body
//	@Tags			GroundObject
//	@Accept			json
//	@Produce		json
//	@Param			ids	body		[]int											true	"Ground Object IDs"
//	@Success		200	{object}	responses.SingleResponse[models.GroundObject]	"Success (returns a blank GroundObject)"
//	@Failure		400	{object}	responses.ErrorResponse							"Invalid JSON or object IDs"
//	@Failure		404	{object}	responses.ErrorResponse							"Objects Not Found"
//	@Failure		500	{object}	responses.ErrorResponse							"Internal Error Deleting Objects"
//	@Router			/groundobject [delete]
func DeleteGroundObjectBatch(c echo.Context) error {
	db, _ := c.Get("db").(*gorm.DB)
	body, _ := io.ReadAll(c.Request().Body)

	var objectIDs []int
	if marshalErr := json.Unmarshal(body, &objectIDs); marshalErr != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Message: "Invalid JSON/ID format",
			Data:    marshalErr.Error()})
	}

	for _, id := range objectIDs {
		if id < 0 {
			return c.JSON(http.StatusBadRequest, responses.ErrorResponse{
				Message: "Invalid ID; Negative ID entered"})
		}
		var objectTBValidated = models.GroundObject{}
		if err := db.First(&objectTBValidated, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, responses.ErrorResponse{
				Message: fmt.Sprintf("Requested object %d does not exist!", id),
				Data:    err.Error()})
		}
	}

	for _, id := range objectIDs {
		dbAction := db.Delete(&models.GroundObject{}, id)
		if err := dbAction.Error; err != nil {
			return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Message: fmt.Sprintf("Error whilst deleting object with id %d", id)})
		}
	}

	return c.JSON(http.StatusOK, responses.SingleResponse[models.GroundObject]{
		Message: "GroundObjects deleted!",
		Model:   models.GroundObject{},
	})
}

// GetAllGroundObjects gets all ground objects in the database
//
//	@Summary		Get all ground objects
//	@Description	Get all ground objects in the database
//	@Tags			GroundObject
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.MultipleResponse[models.GroundObject]	"Success"
//	@Failure		500	{object}	responses.ErrorResponse							"Internal Error Querying GroundObjects"
//	@Router			/groundobjects [get]
func GetAllGroundObjects(c echo.Context) error {
	var objects []models.GroundObject
	db, _ := c.Get("db").(*gorm.DB)

	if err := db.Find(&objects).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Error whilst querying objects!",
			Data:    err.Error()})
	}

	return c.JSON(http.StatusOK, responses.MultipleResponse[models.GroundObject]{
		Message: "GroundObject found!",
		Models:  objects,
	})
}
