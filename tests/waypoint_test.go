package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gcom-backend/configs"
	"gcom-backend/controllers"
	"gcom-backend/models"
	"gcom-backend/responses"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

/*
We can group tests into suites to perform setup and teardown tasks after the suite, test or subtest without having
to constantly call a helper function.
*/

// Here the test suite is defined, we pass in Echo and the DB instance as well
type WaypointTestSuite struct {
	suite.Suite
	e  *echo.Echo
	db *gorm.DB
}

// This is the function which runs all the tests in the suite
func TestRunWaypointSuite(t *testing.T) {
	suite.Run(t, new(WaypointTestSuite))
}

// This function is run when the suite is started
func (s *WaypointTestSuite) SetupSuite() {
	db := configs.ConnectDatabase()
	e := echo.New()

	s.e = e
	s.db = db
}

// This function is run when the suite is finished
func (s *WaypointTestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	err := sqlDB.Close()

	if err != nil {
		fmt.Println("[Teardown] Error closing database connection!")
	}

	err = os.Remove("./database.db")

	if err != nil {
		fmt.Println("[Teardown] Error deleting database!")
	}

}

// This function is run when each test in the suite is finished
func (s *WaypointTestSuite) TearDownTest() {
	s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Waypoint{})
}

// Explanation of difference between assert and require is in this function
func (s *WaypointTestSuite) TestCreateWaypointA() {
	var wp = models.Waypoint{
		ID:          -1,
		Name:        "Alpha",
		Longitude:   123.456,
		Latitude:    123.456,
		Altitude:    100,
		Radius:      10,
		Designation: "land",
		Remarks:     "landing zone",
	}

	var expectedWp = models.Waypoint{
		ID:          1,
		Name:        "Alpha",
		Longitude:   123.456,
		Latitude:    123.456,
		Altitude:    100,
		Radius:      10,
		Designation: "land",
		Remarks:     "landing zone",
	}
	var wpBytes, marshalErr = json.Marshal(&wp)

	if marshalErr != nil {
		panic("Error marshalling test waypoints!")
	}

	c, rec := JSONContextBuilder(s, http.MethodPost, "/waypoint", wpBytes)
	require.NoError(s.T(), controllers.CreateWaypoint(c)) //require stops execution of the test if it fails
	assert.Equal(s.T(), http.StatusOK, rec.Code)          //assert does not stop execution of the test if it fails

	var response responses.WaypointResponse
	require.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &response))

	assert.Equal(s.T(), expectedWp, response.Waypoint)

}

func (s *WaypointTestSuite) TestMalformedJSON() {
	var MalformedJSON = `{
		"namdfe": "Alpha",
		"long": 123.456,
	  	"lat": 123.456,
		"alt": 123.456,
		"radius": 10,
		"designation": "launch",
		"remarks": "launch zone"
	}`

	c, rec := JSONContextBuilder(s, http.MethodPost, "/waypoint", []byte(MalformedJSON))

	if assert.NoError(s.T(), controllers.CreateWaypoint(c)) {
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	}
}

// Example of how to add parameters to the URI is in this function
func (s *WaypointTestSuite) TestEditWaypoint() {
	var wp = models.Waypoint{
		ID:          -1,
		Name:        "Alpha",
		Longitude:   123.456,
		Latitude:    123.456,
		Altitude:    100,
		Radius:      10,
		Designation: "land",
		Remarks:     "landing zone",
	}
	var editJSON = `{"name": "Whiskey"}`
	var expectedWp = models.Waypoint{
		Name:        "Whiskey",
		Longitude:   123.456,
		Latitude:    123.456,
		Altitude:    100,
		Radius:      10,
		Designation: "land",
		Remarks:     "landing zone",
	}

	var wpBytes, marshalErr = json.Marshal(&wp)

	if marshalErr != nil {
		panic("Error marshalling test waypoint!")
	}

	createC, createRec := JSONContextBuilder(s, http.MethodPost, "/waypoint", wpBytes)

	require.NoError(s.T(), controllers.CreateWaypoint(createC))

	var createResp responses.WaypointResponse
	require.NoError(s.T(), json.Unmarshal(createRec.Body.Bytes(), &createResp))
	expectedWp.ID = createResp.Waypoint.ID

	//In order to add parameters, we must specify them this way
	c, rec := JSONContextBuilder(s, http.MethodPatch, "/waypoint", []byte(editJSON))
	c.SetPath("/waypoint/:waypointId")
	c.SetParamNames("waypointId")
	c.SetParamValues(strconv.Itoa(expectedWp.ID))
	require.NoError(s.T(), controllers.EditWaypoint(c))
	assert.Equal(s.T(), http.StatusOK, rec.Code)

	var response responses.WaypointResponse
	require.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &response))

	assert.Equal(s.T(), expectedWp, response.Waypoint)
}

func (s *WaypointTestSuite) TestEditNonExistentWaypoint() {
	var editJSON = `{"name": "Whiskey"}`
	c, rec := JSONContextBuilder(s, http.MethodPut, "/waypoint", []byte(editJSON))
	c.SetPath("/waypoint/:waypointId")
	c.SetParamNames("waypointId")
	c.SetParamValues("1")

	require.NoError(s.T(), controllers.EditWaypoint(c))
	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
}

func (s *WaypointTestSuite) TestQueryWaypoint() {
	var wp = models.Waypoint{
		ID:          -1,
		Name:        "Alpha",
		Longitude:   123.456,
		Latitude:    123.456,
		Altitude:    100,
		Radius:      10,
		Designation: "land",
		Remarks:     "landing zone",
	}

	var expectedWp = models.Waypoint{
		Name:        "Alpha",
		Longitude:   123.456,
		Latitude:    123.456,
		Altitude:    100,
		Radius:      10,
		Designation: "land",
		Remarks:     "landing zone",
	}
	var wpBytes, marshalErr = json.Marshal(&wp)

	if marshalErr != nil {
		panic("Error marshalling test waypoints!")
	}

	createC, createRec := JSONContextBuilder(s, http.MethodPost, "/waypoint", wpBytes)

	require.NoError(s.T(), controllers.CreateWaypoint(createC))

	var createResp responses.WaypointResponse
	require.NoError(s.T(), json.Unmarshal(createRec.Body.Bytes(), &createResp))
	expectedWp.ID = createResp.Waypoint.ID

	c, rec := BlankContextBuilder(s, http.MethodGet, "/waypoint")
	c.SetPath("/waypoint/:waypointId")
	c.SetParamNames("waypointId")
	c.SetParamValues(strconv.Itoa(expectedWp.ID))

	var response responses.WaypointResponse
	require.NoError(s.T(), controllers.GetWaypoint(c))
	require.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &response))
	assert.Equal(s.T(), http.StatusOK, rec.Code)

	assert.Equal(s.T(), expectedWp, response.Waypoint)
}

func (s *WaypointTestSuite) TestQueryNonExistentWaypoint() {
	c, rec := BlankContextBuilder(s, http.MethodGet, "/waypoint")
	c.SetPath("/waypoint/:waypointId")
	c.SetParamNames("waypointId")
	c.SetParamValues("1")

	require.NoError(s.T(), controllers.GetWaypoint(c))
	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
}

func (s *WaypointTestSuite) TestDeleteWaypoint() {
	var wp = models.Waypoint{
		ID:          -1,
		Name:        "Alpha",
		Longitude:   123.456,
		Latitude:    123.456,
		Altitude:    100,
		Radius:      10,
		Designation: "land",
		Remarks:     "landing zone",
	}

	var wpBytes, marshalErr = json.Marshal(&wp)

	if marshalErr != nil {
		panic("Error marshalling test waypoints!")
	}

	createC, createRec := JSONContextBuilder(s, http.MethodPost, "/waypoint", wpBytes)

	require.NoError(s.T(), controllers.CreateWaypoint(createC))

	var createResp responses.WaypointResponse
	require.NoError(s.T(), json.Unmarshal(createRec.Body.Bytes(), &createResp))

	c, rec := BlankContextBuilder(s, http.MethodDelete, "/waypoint")
	c.SetPath("/waypoint/:waypointId")
	c.SetParamNames("waypointId")
	c.SetParamValues(strconv.Itoa(createResp.Waypoint.ID))

	require.NoError(s.T(), controllers.GetWaypoint(c))
	assert.Equal(s.T(), http.StatusOK, rec.Code)
}

func (s *WaypointTestSuite) TestDeleteNonExistentWaypoint() {
	c, rec := BlankContextBuilder(s, http.MethodDelete, "/waypoint")
	c.SetPath("/waypoint/:waypointId")
	c.SetParamNames("waypointId")
	c.SetParamValues("1")

	require.NoError(s.T(), controllers.GetWaypoint(c))
	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
}

func (s *WaypointTestSuite) TestQueryAllWaypoints() {
	var wps = []models.Waypoint{
		{
			ID:          -1,
			Name:        "Alpha",
			Longitude:   123.456,
			Latitude:    123.456,
			Altitude:    100,
			Radius:      10,
			Designation: "land",
			Remarks:     "landing zone",
		},
		{
			ID:          -1,
			Name:        "Beta",
			Longitude:   123.456,
			Latitude:    123.456,
			Altitude:    100,
			Radius:      10,
			Designation: "takeoff",
			Remarks:     "takeoff zone",
		},
	}

	var expectedWps []models.Waypoint

	for _, wp := range wps {
		var wpBytes, marshalErr = json.Marshal(&wp)

		if marshalErr != nil {
			panic("Error marshalling test waypoints!")
		}

		createC, createRec := JSONContextBuilder(s, http.MethodPost, "/waypoint", wpBytes)
		require.NoError(s.T(), controllers.CreateWaypoint(createC))

		var createResp responses.WaypointResponse
		require.NoError(s.T(), json.Unmarshal(createRec.Body.Bytes(), &createResp))

		assert.Equal(s.T(), http.StatusOK, createRec.Code)
		expectedWps = append(expectedWps, createResp.Waypoint)
	}

	c, rec := BlankContextBuilder(s, http.MethodGet, "/waypoints")

	var response responses.WaypointsResponse
	require.NoError(s.T(), controllers.GetAllWaypoints(c))
	require.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &response))
	assert.Equal(s.T(), http.StatusOK, rec.Code)

	assert.Equal(s.T(), expectedWps, response.Waypoints)

}

// These functions help build contexsts that are repeatedly used
func JSONContextBuilder(s *WaypointTestSuite, method string, uri string, m []byte) (echo.Context, *httptest.ResponseRecorder) {
	var req = httptest.NewRequest(method, uri, bytes.NewReader(m))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	var rec = httptest.NewRecorder()
	var c = s.e.NewContext(req, rec)
	//It is important to set the "db" field in the context as controllers expect that to be there
	c.Set("db", s.db)

	return c, rec
}

func BlankContextBuilder(s *WaypointTestSuite, method string, uri string) (echo.Context, *httptest.ResponseRecorder) {
	var req = httptest.NewRequest(method, uri, nil)
	var rec = httptest.NewRecorder()
	var c = s.e.NewContext(req, rec)
	c.Set("db", s.db)

	return c, rec
}
