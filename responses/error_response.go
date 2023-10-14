package responses

// ErrorResponse describes a JSON response for any error
//
// @Description JSON response for any error
type ErrorResponse struct {
	Message string `json:"message" example:"Sample error message"`
	Data    string `json:"data,omitempty"`
}
