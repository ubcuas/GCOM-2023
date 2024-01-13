package configs

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

// type MissionPlannerWrapper interface {
// 	GetQueue()
// 	PostQueue(waypoints []models.Waypoint)
// 	GetStatus()
// 	Lock()
// 	Unlock()
// 	Takeoff(altitude float64)
// 	Land()
// 	ReturnHome()
// 	SetHome(waypoint models.Waypoint)
// }

type MissionPlanner struct {
	url string
}

// ConnectMissionPlanner creates a new instance of MissionPlanner - this should only be in main.go
func ConnectMissionPlanner(url string) (*MissionPlanner, error) {
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

func (mp MissionPlanner) GetQueue() []models.Waypoint {
	resp := genericGet(mp.url + "/queue")
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

func (mp MissionPlanner) GetStatus() *http.Response {
	return genericGet(mp.url + "/status")
}

func (mp MissionPlanner) ReturnHome() bool {
	resp := genericGet(mp.url + "/rtl")
	return resp.StatusCode == http.StatusOK
}

func (mp MissionPlanner) Land() bool {
	resp := genericGet(mp.url + "/land")
	return resp.StatusCode == http.StatusOK
}

func (mp MissionPlanner) Lock() bool {
	resp := genericGet(mp.url + "/lock")
	return resp.StatusCode == http.StatusOK
}

func (mp MissionPlanner) Unlock() bool {
	resp := genericGet(mp.url + "/unlock")
	return resp.StatusCode == http.StatusOK
}

func (mp MissionPlanner) SetQueue(waypoints []models.Waypoint) bool {
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

	resp := genericPost(mp.url+"/queue", json)

	return resp.StatusCode == http.StatusOK
}

func (mp MissionPlanner) Takeoff(alt float64) bool {
	json, err := json.Marshal(map[string]float64{
		"altitude": alt,
	})

	if err != nil {
		log.Fatal("[MP Functions] Error marshalling altitude for takeoff")
	}

	resp := genericPost(os.Getenv("MP_URL"+"/takeoff"), json)

	return resp.StatusCode == http.StatusOK
}

func (mp MissionPlanner) SetHome(waypoint models.Waypoint) bool {
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

	resp := genericPost(mp.url+"/home", json)

	return resp.StatusCode == http.StatusOK
}

func (mp MissionPlanner) SetFlightMode(mode string, drone string, altStandard string) bool {
	json, err := json.Marshal(map[string]string{
		"flight_mode":       mode,
		"drone_type":        drone,
		"altitude_standard": altStandard,
	})

	if err != nil {
		log.Fatal("[MP Functions] Error marshalling flight modes")
	}

	resp := genericPost(mp.url+"/flightmode", json)

	return resp.StatusCode == http.StatusOK
}
