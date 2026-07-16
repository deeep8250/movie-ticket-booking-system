package config

import (
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type dbClient struct {
	PostgresClient *sqlx.DB
	RedisClient    *redis.Client
}

var DBClients dbClient

type Config struct {
	Port string
}

func Load() Config {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return Config{
		Port: ":" + port,
	}
}
