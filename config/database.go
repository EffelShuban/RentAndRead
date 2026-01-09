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
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable prefer_simple_protocol=true",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.Rental{},
		&models.Payment{},
		&models.TopUp{},
	)
	if err != nil {
		log.Fatalf("error migrating database: %v", err)
	}
}