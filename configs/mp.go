package configs

import (
	"net/http"
	"strconv"

	"bytes"
	"encoding/json"
	"gcom-backend/models"
	"io"
	"log"
	"time"
)

type MissionPlanner struct {
	url  string
	lock bool
}

// ConnectMissionPlanner creates a new instance of MissionPlanner - this should only be in main.go
func ConnectMissionPlanner(url string) (*MissionPlanner, error) {
	return &MissionPlanner{
		url:  url,
		lock: false,
	}, nil
	//Assumes method of checking if MP is alive
	// _, err := http.Get(url)
	// if err != nil {
	// 	return &MissionPlanner{
	// 		url: url,
	// 	}, nil
	// } else {
	// 	return nil, errors.New("missionplanner unreachable at provided url")
	// }
}

type mpWaypoint struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Altitude  float64 `json:"altitude"`
}

type mpDrone struct {
	Velocity       float64 `json:"airspeed"`
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
	Altitude       float64 `json:"altitude"`
	Heading        float64 `json:"heading"`
	BatteryVoltage float64 `json:"batteryvoltage"`
}

func genericGet(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		println(url)
		println(err.Error())
		log.Fatal("[MP Functions] Failed running GET")
	}

	return resp
}

func genericPost(url string, json []byte) *http.Response {
	jsonBody := bytes.NewBuffer(json)
	resp, err := http.Post(url, "application/json", jsonBody)
	if err != nil {
		println(url)
		println(jsonBody)
		println(err.Error())
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
			ID:        -1,
			Name:      mpwp.Name,
			Latitude:  mpwp.Latitude,
			Longitude: mpwp.Longitude,
			Altitude:  mpwp.Altitude,
		}

		ans = append(ans, wp)
	}

	return ans
}

func (mp MissionPlanner) GetStatus() models.Drone {
	resp := genericGet(mp.url + "/status")
	var respDrone mpDrone
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("[MP Functions] Error decoding /status response")
	}

	if err := json.Unmarshal(body, &respDrone); err != nil {
		log.Fatal("[MP Functions] Error unmarshalling MP Status")
	}

	var ans = models.Drone{
		Timestamp:      time.Now().Unix(),
		Latitude:       respDrone.Latitude,
		Longitude:      respDrone.Longitude,
		Altitude:       respDrone.Altitude,
		VerticalSpeed:  0.0,
		Speed:          respDrone.Velocity,
		Heading:        respDrone.Heading,
		BatteryVoltage: respDrone.BatteryVoltage,
	}

	return ans
}

func (mp MissionPlanner) ReturnHome(alt float64) bool {
	json, err := json.Marshal(map[string]float64{
		"altitude": alt,
	})

	if err != nil {
		log.Fatal("[MP Functions] Error marshalling altitude for rtl")
	}

	resp := genericPost(mp.url+"/rtl", json)

	return resp.StatusCode == http.StatusOK
}

func (mp MissionPlanner) Land() bool {
	resp := genericGet(mp.url + "/land")
	return resp.StatusCode == http.StatusOK
}

func (mp MissionPlanner) Lock() bool {
	resp := genericGet(mp.url + "/lock")
	mp.lock = true
	return resp.StatusCode == http.StatusOK
}

func (mp MissionPlanner) Unlock() bool {
	resp := genericGet(mp.url + "/unlock")
	mp.lock = false
	return resp.StatusCode == http.StatusOK
}

func (mp MissionPlanner) SetQueue(waypoints []models.Waypoint) bool {
	var mpArr []mpWaypoint
	for _, wp := range waypoints {
		mpwp := mpWaypoint{
			ID:        strconv.Itoa(wp.ID),
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

	resp := genericPost(mp.url+"/takeoff", json)

	return resp.StatusCode == http.StatusOK
}

func (mp MissionPlanner) SetHome(waypoint models.Waypoint) bool {
	mpwp := mpWaypoint{
		ID:        strconv.Itoa(waypoint.ID),
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
