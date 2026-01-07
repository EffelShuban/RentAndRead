package middleware

import (
	"rent-and-read/config"

	echojwt "github.com/labstack/echo-jwt/v4"
)

func JWTMiddleware() echojwt.Config{
	return echojwt.Config{
		SigningKey: config.JwtSecret,
		ContextKey: "user",
	}
}