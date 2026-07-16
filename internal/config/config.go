package config

import (
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

var PostgresClient *sqlx.DB
var RedisClient *redis.Client

type Config struct {
	Port string
}

func Load() Config {
	port := os.Getenv("PORT")

	if port == "" {
		port = ":8080"
	}

	return Config{
		Port: ":" + port,
	}
}
