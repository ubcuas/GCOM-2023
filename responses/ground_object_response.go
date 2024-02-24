package responses

import (
	"gcom-backend/models"
)

// GroundObjectResponse describes a JSON response with a single Ground Object
//
// @Description Describes a response a single ground object
type GroundObjectResponse struct {
	Message  string          `json:"message" example:"Sample success message"`
	GroundObject models.GroundObject `json:"ground_object"`
}
