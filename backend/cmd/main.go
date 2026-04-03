package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/handlers"
	"github.com/ilushew/udmurtia-trip/backend/internal/migrations"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
	"github.com/ilushew/udmurtia-trip/backend/internal/services"
	"github.com/ilushew/udmurtia-trip/backend/pkg/migrator"
	"github.com/ilushew/udmurtia-trip/backend/pkg/postgres"
	redisPkg "github.com/ilushew/udmurtia-trip/backend/pkg/redis"
)

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/base.html", "templates/index.html")
	r.AddFromFiles("register", "templates/base.html", "templates/auth/register.html")
	r.AddFromFiles("verify", "templates/partials/verify-code-form.html")
	r.AddFromFiles("auth-error", "templates/partials/auth-error.html")
	r.AddFromFiles("generate", "templates/partials/route-result.html")
	return r
}

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

	// Инициализация email сервиса
	emailCfg := services.EmailConfig{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     465,
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		From:     os.Getenv("EMAIL_FROM"),
	}
	emailSvc, err := services.NewEmailService(emailCfg)
	if err != nil {
		log.Fatalf("Failed to create email service: %v", err)
	}
	defer emailSvc.Close()

	redisAddr := "redis:6379"
	redisClient := redisPkg.NewClient(redisAddr)

	err = redisPkg.Ping(ctx, redisClient)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	codeService := services.NewCodeService(redisClient)
	userRepo := repository.NewUserRepository(pool)
	authHandler := handlers.NewAuthHandler(userRepo, emailSvc, codeService)
	r := gin.Default()

	// Инициализация сессий
	secret := "your-super-secret-key-min-32-characters-long"
	if secret == "" {
		// Для разработки: сгенерируйте ключ один раз и сохраните в .env
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

	// 👇 Регистрируем middleware (ОБЯЗАТЕЛЬНО до маршрутов!)
	r.Use(sessions.Sessions("udmurtia_trip", store))

	// загрузка шаблонов
	r.HTMLRender = createMyRender()

	r.Static("/static", "static")

	tripHandler := handlers.NewTripHandler()
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{"title": "Home"})
	})
	r.POST("/generate", tripHandler.GenerateTrip)
	auth := r.Group("/auth")
	{
		auth.GET("/register", authHandler.ShowRegisterPage)
		auth.POST("/register", authHandler.Register)
		auth.POST("/verify", authHandler.VerifyCode)
	}
	r.Run(":8080")
}
