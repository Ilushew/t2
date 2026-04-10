package handlers

import (
	"fmt"
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
	app := models.Application{
		RouteName: c.PostForm("route_name"),
		Email:     c.PostForm("email"),
		Comment:   c.PostForm("comment"),
	}

	// Валидация
	if app.Email == "" {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusBadRequest, `<div class="status-error">Email обязателен</div>`)
		return
	}

	// Получаем email админа из конфига
	adminEmail := config.Get("APPLICATION_ADMIN_EMAIL", "")
	if adminEmail == "" {
		log.Printf("APPLICATION_ADMIN_EMAIL не настроен")
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusInternalServerError, `<div class="status-error">Сервис заявок временно недоступен</div>`)
		return
	}

	// Отправляем письмо админу
	routeName := app.RouteName
	if routeName == "" {
		routeName = "Не указан"
	}
	if err := h.emailSvc.SendApplicationToAdmin(adminEmail, routeName, app.Email, app.Comment); err != nil {
		log.Printf("Ошибка отправки заявки админу: %v", err)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusInternalServerError, `<div class="status-error">Ошибка отправки заявки</div>`)
		return
	}

	// Отправляем подтверждение клиенту
	if err := h.emailSvc.SendApplicationConfirmation(app.Email, routeName); err != nil {
		log.Printf("Ошибка отправки подтверждения клиенту: %v", err)
		// Не фейлим запрос, т.к. админ уже получил заявку
	}

	log.Printf("Заявка отправлена от %s (маршрут: %s)", app.Email, routeName)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, fmt.Sprintf(`<div class="status-success">Заявка успешно отправлена! Мы свяжемся с вами в ближайшее время.</div>`))
}
