package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/ilushew/udmurtia-trip/backend/internal/migrations"
	"github.com/ilushew/udmurtia-trip/backend/internal/services"
	"github.com/ilushew/udmurtia-trip/backend/pkg/config"
	"github.com/ilushew/udmurtia-trip/backend/pkg/migrator"
	"github.com/ilushew/udmurtia-trip/backend/pkg/postgres"
	redisPkg "github.com/ilushew/udmurtia-trip/backend/pkg/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func initPostgres(ctx context.Context) (*pgxpool.Pool, *sql.DB, error) {
	pool, err := postgres.NewPool(ctx, postgres.Config{
		Host:     config.Get("DB_HOST", "localhost"),
		Port:     config.Get("DB_PORT", "5432"),
		User:     config.Get("DB_USER", "postgres"),
		Password: config.Get("DB_PASSWORD", "postgres"),
		Database: config.Get("DB_NAME", "udmurtia-trip"),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create pool: %w", err)
	}
	stdDB := postgres.NewStdDB(pool)
	return pool, stdDB, nil
}

func runMigrations(stdDB *sql.DB) error {
	migr, err := migrator.EmbedMigrations(stdDB, migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	if err := migr.Up(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Migrations applied successfully")
	return nil
}

func initEmailService() (*services.EmailService, error) {
	return services.NewEmailService(services.EmailConfig{
		Host:     config.MustGet("EMAIL_HOST"),
		Port:     config.MustGet("EMAIL_PORT"),
		Username: config.MustGet("EMAIL_USERNAME"),
		Password: config.MustGet("EMAIL_PASSWORD"),
		From:     config.MustGet("EMAIL_FROM"),
	})
}

func initRedis(ctx context.Context) (*redis.Client, *services.CodeService, error) {
	client := redisPkg.NewClient(
		config.MustGet("REDIS_ADDR"),
		config.MustGet("REDIS_PASSWORD"),
		config.MustGet("REDIS_DB"),
	)
	if err := redisPkg.Ping(ctx, client); err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return client, services.NewCodeService(client), nil
}
