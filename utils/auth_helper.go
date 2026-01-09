package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetUserID(c echo.Context) uint{
	user, ok := c.Get("user").(*jwt.Token)
	if !ok{
		return 0
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok{
		return 0
	}

	return uint(claims["userID"].(float64))
}