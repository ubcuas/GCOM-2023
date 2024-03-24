package responses

// WaypointsResponse describes a JSON response with multiple Waypoints
//
// @Description Describes a response with multiple waypoints
type MultipleResponse[T any] struct {
	Message   string            `json:"message" example:"Sample success message"`
	Models []T `json:"waypoints"`
}
