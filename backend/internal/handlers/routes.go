package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/services"
)

type TripHandler struct {
	mlService *services.MockMLService
}

// Конструктор хендлера
func NewTripHandler() *TripHandler {
	return &TripHandler{
		mlService: services.NewMockMLService(),
	}
}

func (h *TripHandler) GenerateTrip(c *gin.Context) {
	// 1. Получаем предпочтение из формы
	pref := c.PostForm("trip_type")

	// 2. Вызываем заглушку
	places, _, err := h.mlService.GetRecommendations(pref)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{"msg": "Ошибка генерации"})
		return
	}

	// 3. Отдаем HTML-фрагмент с результатами
	c.HTML(http.StatusOK, "route-result.html", gin.H{
		"places": places,
	})
}
