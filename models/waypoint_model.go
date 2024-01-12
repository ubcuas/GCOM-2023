package models

// Designation describes a Waypoint's designation
//
// @Description Describes a special purpose for a Waypoint
type Designation string

const (
	Launch   Designation = "launch"
	Land     Designation = "land"
	Obstacle Designation = "obstacle"
	// Payload  Designation = "payload"
)

// Waypoint describes a location
//
// @Description describes a location in GCOM
type Waypoint struct {
	//To create a waypoint, ID of "-1" must be passed
	ID        int     `json:"id,string" validate:"required" gorm:"primaryKey" example:"1" extensions:"x-order=1"`
	Name      string  `json:"name" validate:"required" example:"Alpha" extensions:"x-order=2"`
	Latitude  float64 `json:"lat" validate:"required" example:"49.267941" extensions:"x-order=3"`
	Longitude float64 `json:"long" validate:"required" example:"-123.247360" extensions:"x-order=4"`
	Altitude  float64 `json:"alt" validate:"required" example:"100.00" extensions:"x-order=5"`
	//Radius around waypoint where it is considered flown over
	Radius float64 `json:"radius,omitempty" example:"10.0" extensions:"x-order=6"`
	//Designation of waypoint, none by default
	Designation Designation `json:"designation,omitempty" example:"land" extensions:"x-order=7"`
	Remarks     string      `json:"remarks,omitempty" example:"Task 1 Landing Zone" extensions:"x-order=8"`
}
