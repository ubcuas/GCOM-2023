package models

type GroundObj struct {
	//To create a ground object, ID of "-1" must be passed
	ID           int     `json:"id,string" validate:"required" gorm:"primaryKey" example:"-1" extensions:"x-order=1"`
	Object_Type  string  `json:"obj_type" validate:"required" example:"standard" extensions:"x-order=2"`
	Shape        string  `json:"shape" validate:"required" example:"triangle" extensions:"x-order=3"`
	Color        string  `json:"color" validate:"required" example:"blue" extensions:"x-order=4"`
	Text         string  `json:"text" validate:"required" example:"A" extensions:"x-order=5"`
	Text_Color   string  `json:"text_color" validate:"required" example:"green" extensions:"x-order=6"`
	Longitude    float64 `json:"long" validate:"required" example:"-123.45" extensions:"x-order=7"`
	Latitude     float64 `json:"lat" validate:"required" example:"123.45" extensions:"x-order=8"`
}
