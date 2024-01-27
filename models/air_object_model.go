package models

// AirObject describes the information from remoteID of other present drones.
//
// @Description describes the information from remoteID of other present drones.
type AirObject struct {
	ID            string  `json:"id" gorm:"primaryKey" validate:"required" example:"FIN87astrdge12k8"`
	Longitude     float64 `json:"longitude" validate:"required" example:"49.260605"`
	Latitude      float64 `json:"latitude" validate:"required" example:"-123.245995"`
	Altitude      float64 `json:"altitude" validate:"required" example:"75"`
	VerticalSpeed float64 `json:"v_speed" validate:"required" example:"-1.2"`
	Speed         float64 `json:"speed" validate:"required" example:"0.56"`
	Heading       float64 `json:"heading" validate:"required" example:"12"`
	Timestamp     int     `json:"timestamp" validate:"required" example:"1700905713"`
}
