package responses

import (
	"gcom-backend/models"
)

// PayloadResponse describes a JSON response with a single Payload
//
// @Description Describes a response a single payload
type PayloadResponse struct {
	Message		string          `json:"message" example:"Sample success message"`
	Payload 	models.Payload	`json:"payload"`
}
