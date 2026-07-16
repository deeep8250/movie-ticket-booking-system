package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/deeep8250/movie-ticket-booking-system/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func DBinit() {
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"

	var err error
	for range 5 {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		config.DBClients.PostgresClient, err = sqlx.ConnectContext(ctx, "postgres", dsn)
		cancel()
		if err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatal("error while connecting with database", err.Error())
	}

	log.Println("database connection successfull")
}
