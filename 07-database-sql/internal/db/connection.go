package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/qwenhn/gin-restful-api/07-database-sql/internal/config"
)

var DB *sql.DB

func InitDB() error {
	connStr := config.NewConfig().DNS()

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("unable to use data source name", err)
	}

	DB.SetMaxIdleConns(3)                   // Maximum number of idle connections
	DB.SetMaxOpenConns(30)                  // Maximum number of open connections
	DB.SetConnMaxLifetime(30 * time.Minute) // Close connections after 30 minutes
	DB.SetConnMaxIdleTime(5 * time.Minute)  // Close idle connections after 5 minutes

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		DB.Close()
		return fmt.Errorf("DB ping error: %w", err)
	}

	log.Println("Connected")

	return nil
}
