package routes

import (
	"golang-user-authentication/controllers"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	e := echo.New()

	auth := e.Group("/api", echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	auth.PUT("/users", controllers.UpdateUser)
	auth.GET("/users", controllers.GetUser)
	auth.DELETE("/users", controllers.LogoutUser)

	e.POST("/api/users/register", controllers.CreateUser)
	e.POST("/api/users/login", controllers.Login)

	return e
}
