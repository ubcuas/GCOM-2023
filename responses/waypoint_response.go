package responses

import (
	"gcom-backend/models"
)

// WaypointResponse describes a JSON response with a single Waypoint
//
// @Description Describes a response a single waypoint
type WaypointResponse struct {
	Message  string          `json:"message" example:"Sample success message"`
	Waypoint models.Waypoint `json:"waypoint"`
}
