package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/handlers"
	"github.com/ilushew/udmurtia-trip/backend/internal/migrations"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
	"github.com/ilushew/udmurtia-trip/backend/internal/services"
	"github.com/ilushew/udmurtia-trip/backend/pkg/config"
	"github.com/ilushew/udmurtia-trip/backend/pkg/migrator"
	"github.com/ilushew/udmurtia-trip/backend/pkg/postgres"
	redisPkg "github.com/ilushew/udmurtia-trip/backend/pkg/redis"
)

// загружаем шаблоны страниц
func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	// загрузка цельных страниц
	r.AddFromFiles("index", "templates/base.html", "templates/index.html")
	r.AddFromFiles("register", "templates/base.html", "templates/auth/register.html")

	// загрузка частей страниц
	r.AddFromFiles("verify", "templates/partials/verify-code-form.html")
	r.AddFromFiles("auth-error", "templates/partials/auth-error.html")
	r.AddFromFiles("generate", "templates/partials/route-result.html")
	
	return r
}

func main() {
	// загружаем файл конфигураций
	if err := config.LoadDotEnv(".env"); err != nil {
		fmt.Printf("loadDotEnv error: %v", err)
	}

	// подключаем постгрес
	ctx := context.Background()
	cfg := postgres.Config{
		Host:     config.Get("DB_HOST", "localhost"),
		Port:     config.Get("DB_PORT", "5432"),
		User:     config.Get("DB_USER", "postgres"),
		Password: config.Get("DB_PASSWORD", "postgres"),
		Database: config.Get("DB_NAME", "udmurtia-trip"),
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

	// инициализация email сервиса
	emailCfg := services.EmailConfig{
		Host:     config.MustGet("EMAIL_HOST"),
		Port:     config.MustGet("EMAIL_PORT"),
		Username: config.MustGet("EMAIL_USERNAME"),
		Password: config.MustGet("EMAIL_PASSWORD"),
		From:     config.MustGet("EMAIL_FROM"),
	}
	emailSvc, err := services.NewEmailService(emailCfg)
	if err != nil {
		log.Fatalf("Failed to create email service: %v", err)
	}
	defer emailSvc.Close()

	// инициализируем Redis
	redisAddr := config.Get("REDIS_ADDR", "redis:6379")
	redisClient := redisPkg.NewClient(redisAddr)

	err = redisPkg.Ping(ctx, redisClient)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	codeService := services.NewCodeService(redisClient)

	// создаем репозиторий для работы с пользователями
	userRepo := repository.NewUserRepository(pool)

	// создаем роутер
	r := gin.Default()

	// инициализация сессий
	secret := config.MustGet("SESSION_SECRET")
	if secret == "" {
		log.Fatal("SESSION_SECRET environment variable is required")
	}

	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 дней
		HttpOnly: true,      // недоступно из JavaScript
		Secure:   false,     // true только для HTTPS в продакшене
		SameSite: http.SameSiteLaxMode,
	})

	// регистрируем middleware
	r.Use(sessions.Sessions("udmurtia_trip", store))

	// загрузка шаблонов
	r.HTMLRender = createMyRender()

	// загрузка статики
	r.Static("/static", "static")

	// инициализация маршрутов
	indexHandler := handlers.NewIndexHandler()
	tripHandler := handlers.NewTripHandler()
	authHandler := handlers.NewAuthHandler(userRepo, emailSvc, codeService)

	r.GET("/", indexHandler.ShowIndexPage)
	r.POST("/generate", tripHandler.GenerateTrip)
	auth := r.Group("/auth")
	{
		auth.GET("/register", authHandler.ShowRegisterPage)
		auth.POST("/register", authHandler.Register)
		auth.POST("/verify", authHandler.VerifyCode)
	}
	r.Run(":8080")
}
