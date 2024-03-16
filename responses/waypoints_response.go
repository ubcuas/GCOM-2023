package responses

import (
	"gcom-backend/models"
)

// WaypointsResponse describes a JSON response with multiple Waypoints
//
// @Description Describes a response with multiple waypoints
type WaypointsResponse struct {
	Message   string            `json:"message" example:"Sample success message"`
	Waypoints []models.Waypoint `json:"waypoints"`
}
