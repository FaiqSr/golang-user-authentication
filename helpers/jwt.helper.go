package helpers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetUserIdFromJwt(c echo.Context) int {
	auth := c.Get("user").(*jwt.Token)

	claims := auth.Claims.(jwt.MapClaims)
	userId := claims["userId"].(float64)
	return int(userId)
}
