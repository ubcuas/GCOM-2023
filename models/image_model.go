package models

type Image struct {
	Timestamp int64  `json:"timestamp" gorm:"primaryKey" validate:"required" example:"1698544781" extensions:"x-order=1"`
	Filename  string `json:"filename" example:"1714898050.png" extensions:"x-order=2"`
}
