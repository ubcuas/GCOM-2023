package models

type Designation string

const (
	Launch   Designation = "launch"
	Land     Designation = "land"
	Obstacle Designation = "obstacle"
	Payload  Designation = "payload"
)

type Waypoint struct {
	ID          int         `json:"id,string" gorm:"primaryKey"`
	Name        string      `json:"name" validate:"required"`
	Longitude   float64     `json:"long" validate:"required"`
	Latitude    float64     `json:"lat" validate:"required"`
	Altitude    float64     `json:"alt" validate:"required"`
	Radius      float64     `json:"radius,omitempty"`
	Designation Designation `json:"designation,omitempty"`
	Remarks     string      `json:"remarks,omitempty"`
}
