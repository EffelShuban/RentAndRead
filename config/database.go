package config

import (
	"fmt"
	"log"
	"os"
	"rent-and-read/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil{
		log.Fatalf("error in loading env for database connection: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error in establishing connection with database")
	}

	DB.AutoMigrate(
		&models.User{}, 
		&models.Book{}, 
		&models.Rental{}, 
		&models.Payment{},
	)
}