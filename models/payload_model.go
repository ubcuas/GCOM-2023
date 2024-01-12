package models

type Payload struct {
	//To create a payload, ID of "-1" must be passed
	ID             int          `json:"id,string" validate:"required" gorm:"primaryKey" example:"1" extensions:"x-order=1"`
	GroundObject  GroundObject  `json:"ground_object" validate:"required" gorm:"-" extensions:"x-order=2"`
}
