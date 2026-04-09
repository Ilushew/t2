package app

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/handlers"
	"github.com/ilushew/udmurtia-trip/backend/internal/middleware"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
	"github.com/ilushew/udmurtia-trip/backend/internal/services"
	"github.com/ilushew/udmurtia-trip/backend/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Deps struct {
	Pool        *pgxpool.Pool
	UserRepo    *repository.UserRepository
	PlaceRepo   *repository.PlaceRepository
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
	// создание хэндлеров
	indexHandler := handlers.NewIndexHandler()
	criteriaHandler := handlers.NewCriteriaHandler(deps.PlaceRepo)
	authHandler := handlers.NewAuthHandler(
		deps.UserRepo,
		deps.EmailSvc,
		deps.CodeService,
	)
	profileHandler := handlers.NewProfileHandler(deps.UserRepo)
	adminHandler := handlers.NewAdminHandler(deps.UserRepo, deps.PlaceRepo)
	applicationHandler := handlers.NewApplicationHandler(deps.EmailSvc)

	// маршруты для админки
	adminGroup := r.Group("/admin", middleware.RequireAdmin(deps.UserRepo))
	{
		adminGroup.GET("/", adminHandler.ShowTables)
		adminGroup.GET("/table/:table", adminHandler.ViewTable)
		adminGroup.GET("/table/:table/create", adminHandler.CreateUserGet)
		adminGroup.POST("/table/:table/create", adminHandler.CreateUserPost)
		adminGroup.GET("/table/:table/edit/:id", adminHandler.EditRowGet)
		adminGroup.POST("/table/:table/edit/:id", adminHandler.EditRowPost)
		adminGroup.POST("/table/:table/delete/:id", adminHandler.DeleteRow)
	}

	// главная страница
	r.GET("/", indexHandler.ShowIndexPage)

	// Маршрут для формы критериев (POST)
	r.POST("/criteria", criteriaHandler.HandleCriteria)

	// Заявка на маршрут
	r.POST("/applications", applicationHandler.SubmitApplication)

	r.GET("/profile", profileHandler.ShowProfilePage)

	auth := r.Group("/auth")
	{
		auth.GET("/register", authHandler.ShowRegisterPage)
		auth.POST("/register", authHandler.Register)
		auth.POST("/verify", authHandler.VerifyCode)
		auth.GET("/logout", authHandler.Logout)
	}
}
