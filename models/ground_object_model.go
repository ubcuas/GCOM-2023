package models

type ObjectType string

const (
	Standard ObjectType = "standard"
	Emergent ObjectType = "emergent"
)

type Color string

const (
	Black  = "black"
	Red    = "red"
	Blue   = "blue"
	Green  = "green"
	Purple = "purple"
	Brown  = "brown"
	Orange = "orange"
)

type Shape string

const (
	Circle        = "circle"
	SemiCircle    = "semicircle"
	QuarterCircle = "quartercircle"
	Triangle      = "triangle"
	Rectangle     = "rectangle"
	Pentagon      = "pentagon"
	Star          = "star"
	Cross         = "cross"
)

// GroundObject describes both emergent and standard targets.
//
// @Description describes targets in GCOM
type GroundObject struct {
	ID        int        `json:"id,string" validate:"required" gorm:"primaryKey" example:"1" extensions:"x-order=1"`
	Type      ObjectType `json:"object_type" validate:"required" example:"emergent" extensions:"x-order=2"`
	Latitude  float64    `json:"lat" validate:"required" example:"49.267941" extensions:"x-order=3"`
	Longitude float64    `json:"long" validate:"required" example:"-123.247360" extensions:"x-order=4"`
	Shape     Shape      `json:shape example:"circle" extensions:"x-order=5"`
	Color     Color      `json:color example:"black" extensions:"x-order=6"`
	Text      string     `json:text example:"A" extensions:"x-order=7`
	TextColor Color      `json:text_color example"white" extensions:"x-order=8"`
}
