package responses

import (
	"gcom-backend/models"
)

type WaypointsResponse struct {
	Status    int               `json:"status"`
	Message   string            `json:"message"`
	Waypoints []models.Waypoint `json:"waypoints"`
}
