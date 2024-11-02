package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB is a global variable for the database connection
var DB *gorm.DB

// ConnectDB initializes the database connection
func ConnectDB() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/latian_rest_api?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}
