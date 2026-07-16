package main

import (
	"github.com/deeep8250/movie-ticket-booking-system/internal/config"
	"github.com/deeep8250/movie-ticket-booking-system/internal/db"
	"github.com/deeep8250/movie-ticket-booking-system/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.DBinit()
	db.RedisInit()
	r := gin.Default()
	routes.Routes(r)
	r.Run(config.Load().Port)
}
