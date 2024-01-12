package models

import "gorm.io/gorm"

/*
	This migrate function creates or updates tables in the database with our
	model definitions. Remember to add your models here, just put a comma after
	the last one
*/

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&Waypoint{}, &Drone{}, &GroundObject{}, &Payload{})
	if err != nil {
		panic(err)
	}
}
