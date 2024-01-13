package util

import (
	"errors"
	"net/http"

	"bytes"
	"encoding/json"
	"gcom-backend/models"
	"io"
	"log"
	"os"
)

type MissionPlannerWrapper interface {
	GetQueue()
	PostQueue(waypoints []models.Waypoint)
	GetStatus()
	Lock()
	Unlock()
	Takeoff(altitude float64)
	Land()
	ReturnHome()
	SetHome(waypoint models.Waypoint)
}

type MissionPlanner struct {
	url string
}

// NewMissionPlanner creates a new instance of MissionPlanner - this should only be in main.go
func NewMissionPlanner(url string) (*MissionPlanner, error) {
	//Assumes method of checking if MP is alive
	res, err := http.Get(url)
	if err != nil && res.StatusCode == 200 {
		return &MissionPlanner{
			url: url,
		}, nil
	} else {
		return nil, errors.New("missionplanner unreachable at provided url")
	}
}

type mpWaypoint struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Altitude  float64 `json:"altitude"`
}

func genericGet(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("[MP Functions] Failed running GET")
	}

	return resp
}

func genericPost(url string, json []byte) *http.Response {
	jsonBody := bytes.NewBuffer(json)
	resp, err := http.Post(url, "application/json", jsonBody)
	if err != nil {
		log.Fatal("[MP Functions] Failed running POST")
	}

	return resp
}

func GetQueue() []models.Waypoint {
	resp := genericGet(os.Getenv("MP_URL") + "/queue")
	var respArr []mpWaypoint
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("[MP Functions] Error decoding /queue response")
	}

	if err := json.Unmarshal(body, &respArr); err != nil {
		log.Fatal("[MP Functions] Error unmarshalling MP Waypoint Array")
	}
	var ans []models.Waypoint

	for _, mpwp := range respArr {
		wp := models.Waypoint{
			ID:        mpwp.ID,
			Name:      mpwp.Name,
			Latitude:  mpwp.Latitude,
			Longitude: mpwp.Longitude,
			Altitude:  mpwp.Altitude,
		}

		ans = append(ans, wp)
	}

	return ans
}

func GetStatus() *http.Response {
	return genericGet(os.Getenv("MP_URL") + "/status")
}

func ReturnHome() bool {
	resp := genericGet(os.Getenv("MP_URL") + "/rtl")
	return resp.StatusCode == http.StatusOK
}

func Land() bool {
	resp := genericGet(os.Getenv("MP_URL") + "/land")
	return resp.StatusCode == http.StatusOK
}

func Lock() bool {
	resp := genericGet(os.Getenv("MP_URL") + "/lock")
	return resp.StatusCode == http.StatusOK
}

func Unlock() bool {
	resp := genericGet(os.Getenv("MP_URL") + "/unlock")
	return resp.StatusCode == http.StatusOK
}

func SetQueue(waypoints []models.Waypoint) bool {
	var mpArr []mpWaypoint
	for _, wp := range waypoints {
		mpwp := mpWaypoint{
			ID:        wp.ID,
			Name:      wp.Name,
			Longitude: wp.Longitude,
			Latitude:  wp.Latitude,
			Altitude:  wp.Altitude,
		}

		mpArr = append(mpArr, mpwp)
	}

	json, err := json.Marshal(mpArr)

	if err != nil {
		log.Fatal("[MP Functions] Error marshalling waypoint queue")
	}

	resp := genericPost(os.Getenv("MP_URL")+"/queue", json)

	return resp.StatusCode == http.StatusOK
}

func Takeoff(alt float64) bool {
	json, err := json.Marshal(map[string]float64{
		"altitude": alt,
	})

	if err != nil {
		log.Fatal("[MP Functions] Error marshalling altitude for takeoff")
	}

	resp := genericPost(os.Getenv("MP_URL"+"/takeoff"), json)

	return resp.StatusCode == http.StatusOK
}

func SetHome(waypoint models.Waypoint) bool {
	mpwp := mpWaypoint{
		ID:        waypoint.ID,
		Name:      waypoint.Name,
		Longitude: waypoint.Longitude,
		Latitude:  waypoint.Latitude,
		Altitude:  waypoint.Altitude,
	}

	json, err := json.Marshal(mpwp)

	if err != nil {
		log.Fatal("[MP Functions] Error marshalling waypoint!")
	}

	resp := genericPost(os.Getenv("MP_URL")+"/home", json)

	return resp.StatusCode == http.StatusOK
}

func SetFlightMode(mode string, drone string, altStandard string) bool {
	json, err := json.Marshal(map[string]string{
		"flight_mode":       mode,
		"drone_type":        drone,
		"altitude_standard": altStandard,
	})

	if err != nil {
		log.Fatal("[MP Functions] Error marshalling flight modes")
	}

	resp := genericPost(os.Getenv("MP_URL")+"/flightmode", json)

	return resp.StatusCode == http.StatusOK
}
