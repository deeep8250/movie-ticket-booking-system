package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/deeep8250/movie-ticket-booking-system/internal/config"
	"github.com/gin-gonic/gin"
)

func CheckHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}

func CheckReady(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := config.DBClients.PostgresClient.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "not ready",
			"database": "postgres",
			"error":    err.Error(),
		})

		return
	}

	if err := config.DBClients.RedisClient.Ping(ctx).Err(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "not ready",
			"database": "redis",
			"error":    err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ready",
		"postgres": "ok",
		"redis":    "ok",
	})

}
