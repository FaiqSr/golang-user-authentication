package database

import (
	"database/sql"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error
var rediss *redis.Client

func MysqlInit() {
	dbConnection := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME")

	db, err = sql.Open("mysql", dbConnection)
	if err != nil {
		panic(err.Error())
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func GetMysqlConnection() *sql.DB {
	return db
}

func RedisInit(host, password string) {
	rediss = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})
}

func GetRedisConnection() *redis.Client {
	return rediss
}
