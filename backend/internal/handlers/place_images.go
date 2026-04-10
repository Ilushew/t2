package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/models"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
)

type PlaceImageHandler struct {
	imageRepo  *repository.PlaceImageRepository
	placeRepo  *repository.PlaceRepository
	uploadDir  string
}

func NewPlaceImageHandler(
	imageRepo *repository.PlaceImageRepository,
	placeRepo *repository.PlaceRepository,
	uploadDir string,
) *PlaceImageHandler {
	return &PlaceImageHandler{
		imageRepo: imageRepo,
		placeRepo: placeRepo,
		uploadDir: uploadDir,
	}
}

// UploadImage обрабатывает загрузку картинки
func (h *PlaceImageHandler) UploadImage(c *gin.Context) {
	placeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID маршрута"})
		return
	}

	// Проверяем что маршрут существует
	places, err := h.placeRepo.GetByIDs(c.Request.Context(), []int{placeID})
	if err != nil || len(places) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Маршрут не найден"})
		return
	}

	// Получаем файл
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не найден"})
		return
	}

	// Создаём директорию если нет
	if err := os.MkdirAll(h.uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания директории"})
		return
	}

	// Формируем имя файла: place_id_уникальное_имя
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("place_%d_%s%s", placeID, file.Filename[:len(file.Filename)-len(ext)], ext)
	dst := filepath.Join(h.uploadDir, filename)

	// Сохраняем файл
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла"})
		return
	}

	// Устанавливаем права 644 для доступа nginx
	if err := os.Chmod(dst, 0644); err != nil {
		log.Printf("Предупреждение: не удалось изменить права файла %s: %v", filename, err)
	}

	// Записываем в БД
	img := &models.PlaceImage{
		PlaceID:  placeID,
		Filename: filename,
	}
	if err := h.imageRepo.Create(c.Request.Context(), img); err != nil {
		os.Remove(dst) // Откатываем файл если БД упала
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка записи в БД"})
		return
	}

	log.Printf("Картинка загружена: %s для place_id=%d", filename, placeID)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Картинка загружена",
		"filename": filename,
	})
}

// DeleteImage удаляет картинку
func (h *PlaceImageHandler) DeleteImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("image_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID картинки"})
		return
	}

	// Получаем информацию о картинке
	// (для простоты предполагаем что знаем имя, но лучше запросить из БД)
	// Тут нужен метод GetByID в репозитории. Добавим его позже если нужно.
	// Сейчас просто удалим запись, а файл оставим (или можно сделать purge).
	
	if err := h.imageRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Картинка удалена"})
}
