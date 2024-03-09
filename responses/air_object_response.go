package responses

import (
	"gcom-backend/models"
)

// WaypointsResponse describes a JSON response with multiple Waypoints
//
// @Description Describes a response with multiple waypoints
type AirObjectsResponse struct {
	Message    string             `json:"message" example:"Sample success message"`
	AirObjects []models.AirObject `json:"air_objects"`
}
