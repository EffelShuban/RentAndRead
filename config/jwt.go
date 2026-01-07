package config

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var JwtSecret []byte

func LoadJwtConfig(){
	err := godotenv.Load()
	if err != nil{
		log.Fatalf("error in loading jwt secret from env: %v", err)
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == ""{
		log.Fatalf("error secret is empty: %v", err)
	}
	JwtSecret = []byte(secret)
}

func GenerateToken(userID uint) (string, error){
	claims := jwt.MapClaims{
		"userID": userID,
		"exp": time.Now().Add(24 * time.Hour * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}