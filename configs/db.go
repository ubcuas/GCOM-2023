package configs

import (
	"gcom-backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	database, err := gorm.Open(sqlite.Open("./db/database.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	models.Migrate(database)

	return database
}
