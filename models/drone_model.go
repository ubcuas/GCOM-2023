package models

// Drone describes the drone being flown
//
// @Description describes the drone being flown
type Drone struct {
	Timestamp     int64     `json:"timestamp" gorm:"primaryKey" validate:"required" example:"1698544781" extensions:"x-order=1"`
	Latitude      float64 `json:"lat" validate:"required" example:"49.267941" extensions:"x-order=2"`
	Longitude     float64 `json:"long" validate:"required" example:"-123.247360" extensions:"x-order=3"`
	Altitude      float64 `json:"alt" validate:"required" example:"100.00" extensions:"x-order=4"`
	VerticalSpeed float64 `json:"v_speed" validate:"required" example:"-1.63" extensions:"x-order=5"`
	Speed         float64 `json:"speed" validate:"required" example:"0.98" extensions:"x-order=6"`
	Heading       float64 `json:"heading" validate:"required" example:"298.12" extensions:"x-order=7"`
	//Payloads TBD
	BatteryVoltage float64 `json:"battery_voltage" validate:"required" example:"2.6" extensions:"x-order=9"`
}
