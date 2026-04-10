package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/models"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
)

type PlaceCommentHandler struct {
	commentRepo *repository.PlaceCommentRepository
	placeRepo   *repository.PlaceRepository
}

func NewPlaceCommentHandler(
	commentRepo *repository.PlaceCommentRepository,
	placeRepo *repository.PlaceRepository,
) *PlaceCommentHandler {
	return &PlaceCommentHandler{
		commentRepo: commentRepo,
		placeRepo:   placeRepo,
	}
}

// GetComments возвращает HTML-фрагмент с комментариями для HTMX
func (h *PlaceCommentHandler) GetComments(c *gin.Context) {
	placeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Неверный ID маршрута")
		return
	}

	// Проверяем что маршрут существует
	places, err := h.placeRepo.GetByIDs(c.Request.Context(), []int{placeID})
	if err != nil || len(places) == 0 {
		c.String(http.StatusNotFound, "Маршрут не найден")
		return
	}

	comments, err := h.commentRepo.GetByPlaceID(c.Request.Context(), placeID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка загрузки комментариев")
		return
	}

	c.HTML(http.StatusOK, "place-comments-fragment", gin.H{
		"PlaceID":  placeID,
		"Comments": comments,
	})
}

// CreateComment обрабатывает отправку формы комментария
func (h *PlaceCommentHandler) CreateComment(c *gin.Context) {
	placeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Неверный ID маршрута")
		return
	}

	// Проверяем что маршрут существует
	places, err := h.placeRepo.GetByIDs(c.Request.Context(), []int{placeID})
	if err != nil || len(places) == 0 {
		c.String(http.StatusNotFound, "Маршрут не найден")
		return
	}

	author := strings.TrimSpace(c.PostForm("author"))
	text := strings.TrimSpace(c.PostForm("text"))

	if author == "" || text == "" {
		c.String(http.StatusBadRequest, "Имя и текст комментария обязательны")
		return
	}

	comment := &models.PlaceComment{
		PlaceID: placeID,
		Author:  author,
		Text:    text,
	}

	if err := h.commentRepo.Create(c.Request.Context(), comment); err != nil {
		c.String(http.StatusInternalServerError, "Ошибка сохранения комментария")
		return
	}

	// Возвращаем обновлённый список комментариев
	comments, err := h.commentRepo.GetByPlaceID(c.Request.Context(), placeID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка загрузки комментариев")
		return
	}

	c.HTML(http.StatusOK, "place-comments-fragment", gin.H{
		"PlaceID":  placeID,
		"Comments": comments,
	})
}
