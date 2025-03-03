package main

import (
	"golang-user-authentication/database"
	"golang-user-authentication/routes"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.MysqlInit()
	database.RedisInit(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD"))

	e := routes.Init()

	e.Logger.Fatal(e.Start(":" + os.Getenv("SERVER_PORT")))

}
