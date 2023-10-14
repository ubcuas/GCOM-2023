package responses

import (
	"gcom-backend/models"
)

type WaypointResponse struct {
	Status   int             `json:"status"`
	Message  string          `json:"message"`
	Waypoint models.Waypoint `json:"waypoint"`
}
