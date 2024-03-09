package configs

import (
	"gcom-backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(test bool) *gorm.DB {
	var db_string = ""
	if test {
		db_string = "database.db"
	} else {
		db_string = "./db/database.db"
	}
	database, err := gorm.Open(sqlite.Open(db_string), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	models.Migrate(database)

	return database
}
