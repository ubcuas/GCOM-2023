package models

type AirObject struct {
	ID            string  `json:"id"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	Altitude      float64 `json:"altitude"`
	VerticalSpeed float64 `json:"v_speed"`
	Speed         float64 `json:"speed"`
	Heading       float64 `json:"heading"`
	Timestamp     int     `json:"timestamp"`
}
