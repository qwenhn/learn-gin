package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/qwenhn/gin-restful-api/09-pgx-sqlc/internal/config"
	"github.com/qwenhn/gin-restful-api/09-pgx-sqlc/internal/db/sqlc"
)

var DB *sqlc.Queries

func InitDB() error {
	connStr := config.NewConfig().DNS()

	conf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("error parsing DB config: %v", err)
	}

	conf.MaxConns = 50
	conf.MinConns = 5
	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute

	ctx := context.Background()

	DBPool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return fmt.Errorf("error creating DB pool: %v", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	DB = sqlc.New(DBPool)

	if err := DBPool.Ping(pingCtx); err != nil {
		DBPool.Close()
		return fmt.Errorf("db ping error: %v", err)
	}

	log.Println("Connected")

	return nil
}
