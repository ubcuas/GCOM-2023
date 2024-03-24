package responses


// SingleResponse describes a JSON response with a single instance of a model
//
// @Description Describes a response a single model instance
type SingleResponse[T any] struct {
	Message  string          `json:"message" example:"Sample success message"`
	Model T `json:"waypoint"`
}
