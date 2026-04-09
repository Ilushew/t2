package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/models"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
)

type CriteriaHandler struct {
	placeRepo *repository.PlaceRepository
}

// Конструктор хендлера
func NewCriteriaHandler(placeRepo *repository.PlaceRepository) *CriteriaHandler {
	return &CriteriaHandler{
		placeRepo: placeRepo,
	}
}

// HandleCriteria обрабатывает отправку формы критериев
func (h *CriteriaHandler) HandleCriteria(c *gin.Context) {
	// Парсим данные из формы
	criteria := models.TripCriteria{
		Duration: c.PostForm("duration"),
		Company:  c.PostForm("company"),
		Budget:   c.PostForm("budget"),
	}

	// Преобразуем has_car
	hasCar := c.PostForm("has_car")
	criteria.HasCar = hasCar == "true"

	// Преобразуем with_pets
	withPets := c.PostForm("with_pets")
	criteria.WithPets = withPets == "true"

	// Собираем интересы (множественный выбор)
	criteria.Interests = c.PostFormArray("interests")

	// Сохраняем JSON в папку (для тестирования)
	if err := saveCriteriaToFile(criteria); err != nil {
		log.Printf("failed to save criteria: %v", err)
		c.HTML(http.StatusOK, "generate", gin.H{
			"message": "Ошибка сохранения критериев",
		})
		return
	}

	// TODO: пока заглушка — в будущем здесь будет вызов ML-сервиса
	// ML вернёт []uuid.UUID — ID мест по убыванию релевантности
	// places, err := h.placeRepo.GetByIDs(ctx, mlIDs)

	c.HTML(http.StatusOK, "generate", gin.H{
		"criteria": criteria,
		"message":  "ML-сервис пока не подключён, данные сохранены в JSON",
	})
}

// saveCriteriaToFile сохраняет критерии в JSON файл в папку criteria_data
func saveCriteriaToFile(criteria models.TripCriteria) error {
	// Создаём папку для данных, если не существует
	dataDir := "criteria_data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Генерируем имя файла с таймстампом
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(dataDir, fmt.Sprintf("criteria_%s.json", timestamp))

	// Маршалим в JSON с отступами
	data, err := json.MarshalIndent(criteria, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal criteria: %w", err)
	}

	// Записываем файл
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("criteria saved to %s", filename)
	return nil
}
