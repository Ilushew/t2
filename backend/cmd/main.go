package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/handlers"
	"github.com/ilushew/udmurtia-trip/backend/internal/migrations"
	"github.com/ilushew/udmurtia-trip/backend/pkg/migrator"
	"github.com/ilushew/udmurtia-trip/backend/pkg/postgres"
)

func main() {
	ctx := context.Background()
	cfg := postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	}
	pool, err := postgres.NewPool(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	stdDB := postgres.NewStdDB(pool)

	// Создаём и запускаем мигратор
	migr, err := migrator.EmbedMigrations(stdDB, migrations.FS, ".")
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}

	if err := migr.Up(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully")

	r := gin.Default()
	// загрузка шаблонов
	r.LoadHTMLFiles(
		"templates/base.html",
		"templates/index.html",
		"templates/partials/route-result.html",
	)
	r.Static("/static", "static")

	tripHandler := handlers.NewTripHandler()
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title": "Главная - Udmurtia AI Route",
		})
	})
	r.POST("/generate", tripHandler.GenerateTrip)
	r.Run(":8080")
}
