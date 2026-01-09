package config

import (
	"os"

	"github.com/joho/godotenv"
)

type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

var Email EmailConfig

func InitEmail() {
	_ = godotenv.Load()

	Email = EmailConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}

	if Email.From == "" {
		Email.From = "no-reply@rentandread.com"
	}
}