package models

type Payload struct {
	//To create a payload, ID of "-1" must be passed
	ID             int           `json:"id,string" validate:"required" gorm:"primaryKey" example:"1" extensions:"x-order=1"`
	Ground_Object  GroundObject  `json:"object" validate:"required" extensions:"x-order=2"`
}
