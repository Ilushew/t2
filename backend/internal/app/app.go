package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
	"github.com/ilushew/udmurtia-trip/backend/internal/services"
	"github.com/ilushew/udmurtia-trip/backend/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	pool        *pgxpool.Pool
	redisClient *redis.Client
	userRepo    *repository.UserRepository
	placeRepo   *repository.PlaceRepository
	emailSvc    *services.EmailService
	codeService *services.CodeService
	router      *gin.Engine
	server      *http.Server
}

func (a *App) Close() {
	if a.redisClient != nil {
		a.redisClient.Close()
	}
	if a.emailSvc != nil {
		a.emailSvc.Close()
	}
	if a.pool != nil {
		a.pool.Close()
	}
}

func New() (*App, error) {
	// загружаем файл конфигураций
	if err := config.LoadDotEnv(".env"); err != nil {
		fmt.Printf("loadDotEnv error: %v", err)
	}

	app := &App{}
	ctx := context.Background()

	// PostgreSQL
	pool, stdDB, err := initPostgres(ctx)
	app.pool = pool
	if err != nil {
		return nil, err
	}

	// миграции
	if err := runMigrations(stdDB); err != nil {
		app.Close()
		return nil, err
	}

	// email сервис
	emailSvc, err := initEmailService()
	app.emailSvc = emailSvc
	if err != nil {
		app.Close()
		return nil, err
	}

	// redis
	redisClient, codeService, err := initRedis(ctx)
	app.redisClient, app.codeService = redisClient, codeService
	if err != nil {
		app.Close()
		return nil, err
	}

	// репозитории
	app.userRepo = repository.NewUserRepository(pool)
	app.placeRepo = repository.NewPlaceRepository(pool)

	// роутер
	deps := Deps{
		Pool:        app.pool,
		UserRepo:    app.userRepo,
		PlaceRepo:   app.placeRepo,
		EmailSvc:    app.emailSvc,
		CodeService: app.codeService,
	}

	app.router = setupRouter(deps)
	app.server = &http.Server{
		Addr:    ":8080",
		Handler: app.router,
	}
	return app, nil
}

func (a *App) Run() {
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	a.Close()
}
