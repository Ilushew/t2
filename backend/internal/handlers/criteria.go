package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

// RouteState хранит состояние пагинации маршрутов
type RouteState struct {
	PlaceIDs []int              `json:"place_ids"`
	Current  int                `json:"current"`
	Criteria models.TripCriteria `json:"criteria"`
}

// HandleCriteria обрабатывает отправку формы критериев
func (h *CriteriaHandler) HandleCriteria(c *gin.Context) {
	// Парсим бюджет как число
	budgetStr := c.PostForm("budget")
	budget, err := strconv.Atoi(budgetStr)
	if err != nil {
		budget = 1 // по умолчанию "Средний"
	}

	// Парсим данные из формы
	criteria := models.TripCriteria{
		Duration: c.PostForm("duration"),
		Company:  c.PostForm("company"),
		Budget:   budget,
		Query:    c.PostForm("preferences"),
	}

	// Преобразуем car (yes/no -> bool)
	hasCar := c.PostForm("car")
	criteria.HasCar = hasCar == "yes"

	// Преобразуем pets (yes/no -> bool)
	withPets := c.PostForm("pets")
	criteria.WithPets = withPets == "yes"

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

	// Вызываем ML-сервис для получения рекомендаций
	mlPlaceIDs, err := callMLService(criteria)
	if err != nil {
		log.Printf("ML service error: %v", err)
		c.HTML(http.StatusOK, "generate", gin.H{
			"message":  "Ошибка получения рекомендаций, попробуйте позже",
			"criteria": criteria,
		})
		return
	}

	// Получаем места по ID от ML
	places, err := h.placeRepo.GetByIDs(c.Request.Context(), mlPlaceIDs)
	if err != nil {
		log.Printf("failed to get places: %v", err)
		c.HTML(http.StatusOK, "generate", gin.H{
			"message":  "Ошибка загрузки маршрутов",
			"criteria": criteria,
		})
		return
	}

	// Сохраняем состояние маршрутов в cookie
	routeState := RouteState{
		PlaceIDs: mlPlaceIDs,
		Current:  0,
		Criteria: criteria,
	}
	setStateCookie(c, routeState)

	// Возвращаем только первый маршрут
	c.HTML(http.StatusOK, "generate", gin.H{
		"criteria":   criteria,
		"place":      places[0],
		"placeNum":   1,
		"total":      len(places),
		"hasMore":    len(places) > 1,
		"yandexKey":  os.Getenv("YANDEX_KEY"),
	})
}

// HandleNextRoute показывает следующий маршрут
func (h *CriteriaHandler) HandleNextRoute(c *gin.Context) {
	// Читаем состояние из cookie
	routeState, err := getStateCookie(c)
	if err != nil {
		c.HTML(http.StatusOK, "generate", gin.H{
			"message": "Сессия маршрутов истекла, начните заново",
		})
		return
	}

	// Увеличиваем счётчик для получения СЛЕДУЮЩЕГО маршрута
	routeState.Current++

	// Проверяем, есть ли ещё маршруты
	if routeState.Current >= len(routeState.PlaceIDs) {
		// Все маршруты просмотрены, обновляем cookie
		setStateCookie(c, routeState)
		c.HTML(http.StatusOK, "generate", gin.H{
			"criteria": routeState.Criteria,
			"message":  "Маршруты закончились! Попробуйте изменить критерии",
		})
		return
	}

	// Получаем место (Current уже увеличен — это следующий маршрут)
	currentPlaceID := routeState.PlaceIDs[routeState.Current]
	places, err := h.placeRepo.GetByIDs(c.Request.Context(), []int{currentPlaceID})
	if err != nil || len(places) == 0 {
		log.Printf("failed to get place %d: %v", currentPlaceID, err)
		c.HTML(http.StatusOK, "generate", gin.H{
			"message": "Ошибка загрузки маршрута",
		})
		return
	}

	// Сохраняем обновлённое состояние
	setStateCookie(c, routeState)

	// Возвращаем место (placeNum = Current + 1, т.к. нумерация с 1)
	c.HTML(http.StatusOK, "generate", gin.H{
		"criteria":  routeState.Criteria,
		"place":     places[0],
		"placeNum":  routeState.Current + 1,
		"total":     len(routeState.PlaceIDs),
		"hasMore":   routeState.Current+1 < len(routeState.PlaceIDs),
		"yandexKey": os.Getenv("YANDEX_KEY"),
	})
}

// HandleRestartRoutes сбрасывает счётчик и показывает маршруты заново
func (h *CriteriaHandler) HandleRestartRoutes(c *gin.Context) {
	// Читаем состояние из cookie
	routeState, err := getStateCookie(c)
	if err != nil {
		c.HTML(http.StatusOK, "generate", gin.H{
			"message": "Сессия маршрутов истекла, начните заново",
		})
		return
	}

	// Сбрасываем счётчик на 0
	routeState.Current = 0
	setStateCookie(c, routeState)

	// Проверяем, есть ли маршруты
	if len(routeState.PlaceIDs) == 0 {
		c.HTML(http.StatusOK, "generate", gin.H{
			"criteria": routeState.Criteria,
			"message":  "Маршруты не найдены, попробуйте изменить критерии",
		})
		return
	}

	// Получаем первое место
	firstPlaceID := routeState.PlaceIDs[0]
	places, err := h.placeRepo.GetByIDs(c.Request.Context(), []int{firstPlaceID})
	if err != nil || len(places) == 0 {
		log.Printf("failed to get place %d: %v", firstPlaceID, err)
		c.HTML(http.StatusOK, "generate", gin.H{
			"message": "Ошибка загрузки маршрута",
		})
		return
	}

	// Возвращаем первый маршрут
	c.HTML(http.StatusOK, "generate", gin.H{
		"criteria": routeState.Criteria,
		"place":    places[0],
		"placeNum": 1,
		"total":    len(routeState.PlaceIDs),
		"hasMore":  len(routeState.PlaceIDs) > 1,
	})
}

// callMLService отправляет критерии в ML-сервис и возвращает список ID мест
func callMLService(criteria models.TripCriteria) ([]int, error) {
	// Формируем query из всех критериев
	queryText := fmt.Sprintf("%s %s %d %v %s %v", criteria.Duration, criteria.Company, criteria.Budget, criteria.HasCar, criteria.Interests, criteria.WithPets)

	// Если есть пользовательские пожелания — добавляем их
	if criteria.Query != "" {
		queryText = fmt.Sprintf("%s %s", criteria.Query, queryText)
	}

	// ML-сервис ожидает плоскую структуру (как в RecommendationRequest)
	mlReq := map[string]any{
		"duration":  criteria.Duration,
		"company":   criteria.Company,
		"has_car":   criteria.HasCar,
		"budget":    criteria.Budget,
		"interests": criteria.Interests,
		"with_pets": criteria.WithPets,
		"query":     queryText,
	}

	body, err := json.Marshal(mlReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Отправляем POST запрос к ml_service
	mlURL := "http://ml_service:8000/predict"
	resp, err := http.Post(mlURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned status: %d", resp.StatusCode)
	}

	// Парсим ответ — массив int ID мест
	var placeIDs []int
	if err := json.NewDecoder(resp.Body).Decode(&placeIDs); err != nil {
		return nil, fmt.Errorf("failed to decode ML response: %w", err)
	}

	log.Printf("ML service returned %d place IDs (query: %s)", len(placeIDs), queryText)
	return placeIDs, nil
}

// setStateCookie сохраняет состояние маршрутов
func setStateCookie(c *gin.Context, state RouteState) {
	data, _ := json.Marshal(state)
	encoded := base64.StdEncoding.EncodeToString(data)
	c.SetCookie("route_state", encoded, 3600, "/", "", false, false) // 1 час
}

// getStateCookie читает состояние маршрутов
func getStateCookie(c *gin.Context) (RouteState, error) {
	encoded, err := c.Cookie("route_state")
	if err != nil {
		return RouteState{}, fmt.Errorf("cookie not found")
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return RouteState{}, fmt.Errorf("decode error: %w", err)
	}

	var state RouteState
	if err := json.Unmarshal(data, &state); err != nil {
		return RouteState{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return state, nil
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
