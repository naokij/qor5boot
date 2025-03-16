package admin

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbParamsString = getEnvWithDefault("DB_PARAMS", "user=docs password=docs dbname=docs sslmode=disable host=localhost port=6532 TimeZone=Asia/Tokyo")

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func ConnectDB() (db *gorm.DB) {
	var err error
	db, err = gorm.Open(postgres.Open(dbParamsString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Logger = db.Logger.LogMode(logger.Info)

	return db
}
