package models

type Image struct {
	Timestamp int64   `json:"timestamp" gorm:"primaryKey" validate:"required" example:"1698544781" extensions:"x-order=1"`
	Filename  string  `json:"filename" example:"1714898050.png" extensions:"x-order=2"`
	Latitude  float64 `json:"lat" example:"49.267941" extensions:"x-order=3"`
	Longitude float64 `json:"long" example:"-123.247360" extensions:"x-order=4"`
	Altitude  float64 `json:"alt" example:"100.00" extensions:"x-order=5"`
	Heading   float64 `json:"heading" example:"298.12" extensions:"x-order=6"`
}
