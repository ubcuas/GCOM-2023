package configs

import (
	"gcom-backend/models"
	"os"

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

	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		os.Mkdir("./db", 0755)
	}

	database, err := gorm.Open(sqlite.Open(db_string), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	models.Migrate(database)

	return database
}
