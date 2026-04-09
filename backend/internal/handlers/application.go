package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/models"
	"github.com/ilushew/udmurtia-trip/backend/internal/services"
	"github.com/ilushew/udmurtia-trip/backend/pkg/config"
)

type ApplicationHandler struct {
	emailSvc *services.EmailService
}

func NewApplicationHandler(emailSvc *services.EmailService) *ApplicationHandler {
	return &ApplicationHandler{
		emailSvc: emailSvc,
	}
}

// SubmitApplication обрабатывает отправку заявки
func (h *ApplicationHandler) SubmitApplication(c *gin.Context) {
	var app models.Application
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных: " + err.Error(),
		})
		return
	}

	// Валидация
	if app.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email обязателен",
		})
		return
	}
	if app.RouteName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Название маршрута обязательно",
		})
		return
	}

	// Получаем email админа из конфига
	adminEmail := config.Get("APPLICATION_ADMIN_EMAIL", "")
	if adminEmail == "" {
		log.Printf("APPLICATION_ADMIN_EMAIL не настроен")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Сервис заявок временно недоступен",
		})
		return
	}

	// Отправляем письмо админу
	if err := h.emailSvc.SendApplicationToAdmin(adminEmail, app.RouteName, app.Email, app.Comment); err != nil {
		log.Printf("Ошибка отправки заявки админу: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка отправки заявки",
		})
		return
	}

	// Отправляем подтверждение клиенту
	if err := h.emailSvc.SendApplicationConfirmation(app.Email, app.RouteName); err != nil {
		log.Printf("Ошибка отправки подтверждения клиенту: %v", err)
		// Не фейлим запрос, т.к. админ уже получил заявку
	}

	log.Printf("Заявка на маршрут '%s' отправлена от %s", app.RouteName, app.Email)

	c.JSON(http.StatusOK, gin.H{
		"message": "Заявка успешно отправлена",
	})
}
