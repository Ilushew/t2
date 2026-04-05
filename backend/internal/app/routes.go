package app

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/handlers"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
	"github.com/ilushew/udmurtia-trip/backend/internal/services"
	"github.com/ilushew/udmurtia-trip/backend/pkg/config"
)

type Deps struct {
	UserRepo    *repository.UserRepository
	EmailSvc    *services.EmailService
	CodeService *services.CodeService
}

func setupRouter(deps Deps) *gin.Engine {
	r := gin.Default()
	setupMiddleware(r)
	setupStaticAndTemplates(r)
	setupHandlers(r, deps)

	return r
}

func setupMiddleware(r *gin.Engine) {
	secret := config.MustGet("SESSION_SECRET")
	if secret == "" {
		panic("SESSION_SECRET env variable is required")
	}

	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 дней
		HttpOnly: true,      // недоступно из JavaScript
		Secure:   false,     // true только для HTTPS в продакшене
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("udmurtia-trip", store))
}

func setupHandlers(r *gin.Engine, deps Deps) {
	indexHandler := handlers.NewIndexHandler()
	tripHandler := handlers.NewTripHandler()
	authHandler := handlers.NewAuthHandler(
		deps.UserRepo,
		deps.EmailSvc,
		deps.CodeService,
	)

	r.GET("/", indexHandler.ShowIndexPage)
	r.POST("/generate", tripHandler.GenerateTrip)

	auth := r.Group("/auth")
	{
		auth.GET("/register", authHandler.ShowRegisterPage)
		auth.POST("/register", authHandler.Register)
		auth.POST("/verify", authHandler.VerifyCode)
		auth.GET("/logout", authHandler.Logout)
	}
}
