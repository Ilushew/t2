package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/handlers"
)

func main() {
	r := gin.Default()

	// Загружаем все шаблоны из папки templates
	r.LoadHTMLFiles(
		"templates/base.html",
		"templates/index.html",
		"templates/partials/route-result.html",
	)
	// Раздаем статику (если понадобится)
	r.Static("/static", "static")

	// Создаем хендлер
	tripHandler := handlers.NewTripHandler()

	// Главный маршрут
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title": "Главная - Udmurtia AI Route",
		})
	})
	r.POST("/generate", tripHandler.GenerateTrip)
	// Запуск сервера на порту 8080
	r.Run(":8080")
}
