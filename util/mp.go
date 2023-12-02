package util

import (
	"errors"
	"net/http"
)

type FlightMode string

const (
	Loiter    FlightMode = "loiter"
	Stabilize FlightMode = "stabilize"
	Auto      FlightMode = "auto"
	Guided    FlightMode = "guided"
)

type DroneType string

const (
	VTOL  DroneType = "vtol"
	Plane DroneType = "plane"
)

type AltitudeType string

const (
	AboveSeaLevel    AltitudeType = "ASL"
	AboveGroundLevel AltitudeType = "AGL"
)

//type MissionPlanner interface {
//	getQueue()
//	setQueue(waypoints []models.Waypoint)
//	getStatus()
//	lock()
//	unlock()
//	takeoff(altitude float64)
//	land()
//	returnHome()
//	setHome(waypoint models.Waypoint)
//}

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
		return nil, errors.New("invalid URL for MissionPlanner")
	}
}
