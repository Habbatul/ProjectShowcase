package config

import (
	"awesomeProject/main/entity"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=testfirstgo_db port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = database
	fmt.Println("Database connected!")
}

func MigrateDatabase() {
	err := DB.AutoMigrate(&entity.Project{}, &entity.Image{}, &entity.Tag{})
	if err != nil {
		return
	}
}
