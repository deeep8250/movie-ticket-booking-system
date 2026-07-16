package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/deeep8250/movie-ticket-booking-system/internal/config"
	"github.com/go-redis/redis/v8"
)

func RedisInit() {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	config.DBClients.RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := config.DBClients.RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("redis connection failed", err.Error())
	}
	log.Println("redis connection is successful")
}
