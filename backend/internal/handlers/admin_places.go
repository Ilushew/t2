package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/models"
)

// viewPlaces — просмотр таблицы places
func (h *AdminHandler) viewPlaces(c *gin.Context) {
	places, err := h.placeRepo.GetAll(c.Request.Context())
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin-error", h.aH(c, gin.H{"message": err.Error()}))
		return
	}

	columns := []struct{ Name, Type string }{
		{"id", "int"},
		{"name", "text"},
		{"price", "text"},
		{"time", "float"},
		{"types_of_movement", "text"},
		{"category", "text"},
		{"is_indoor", "bool"},
		{"with_child", "bool"},
		{"with_pets", "bool"},
	}

	var rows [][]string
	for _, p := range places {
		rows = append(rows, []string{
			fmt.Sprintf("%d", p.ID),
			p.Name,
			p.Price,
			fmt.Sprintf("%.1f", p.Time),
			p.TypesOfMovement,
			p.Category,
			fmt.Sprintf("%v", p.IsIndoor),
			fmt.Sprintf("%v", p.WithChild),
			fmt.Sprintf("%v", p.WithPets),
		})
	}

	c.HTML(http.StatusOK, "admin-view", h.aH(c, gin.H{
		"Title":     "Таблица: places",
		"TableName": "places",
		"Columns":   columns,
		"Rows":      rows,
		"PKName":    "id",
	}))
}

// createPlaceGet — форма добавления места
func (h *AdminHandler) createPlaceGet(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-place-create", h.aH(c, gin.H{
		"Title":     "Добавить место",
		"TableName": "places",
	}))
}

// createPlacePost — создание места
func (h *AdminHandler) createPlacePost(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.HTML(http.StatusBadRequest, "admin-error", h.aH(c, gin.H{"message": "Название обязательно"}))
		return
	}

	price := c.PostForm("price")
	if price == "" {
		price = "mid"
	}

	place := models.Place{
		Name:            name,
		Price:           price,
		TypesOfMovement: c.PostForm("types_of_movement"),
		Category:        c.PostForm("category"),
		Description:     c.PostForm("description"),
		IsIndoor:        c.PostForm("is_indoor") == "on",
		WithChild:       c.PostForm("with_child") == "on",
		WithPets:        c.PostForm("with_pets") == "on",
	}

	fmt.Sscanf(c.PostForm("time"), "%f", &place.Time)
	fmt.Sscanf(c.PostForm("lat_start"), "%f", &place.LatStart)
	fmt.Sscanf(c.PostForm("lon_start"), "%f", &place.LonStart)
	fmt.Sscanf(c.PostForm("lat_end"), "%f", &place.LatEnd)
	fmt.Sscanf(c.PostForm("lon_end"), "%f", &place.LonEnd)

	err := h.placeRepo.InsertPlace(c.Request.Context(), &place)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin-error", h.aH(c, gin.H{"message": err.Error()}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/table/places")
}

// editPlaceGet — форма редактирования места
func (h *AdminHandler) editPlaceGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin-error", h.aH(c, gin.H{"message": "Неверный ID"}))
		return
	}

	place, err := h.findPlaceByID(c, id)
	if err != nil {
		c.HTML(http.StatusNotFound, "admin-error", h.aH(c, gin.H{"message": err.Error()}))
		return
	}

	c.HTML(http.StatusOK, "admin-edit", h.aH(c, gin.H{
		"Title":     "Редактирование: places",
		"TableName": "places",
		"Columns": []string{
			"name", "price", "time", "types_of_movement", "category",
			"is_indoor", "with_child", "with_pets", "description",
		},
		"RowData": map[string]any{
			"name":              place.Name,
			"price":             place.Price,
			"time":              fmt.Sprintf("%.1f", place.Time),
			"types_of_movement": place.TypesOfMovement,
			"category":          place.Category,
			"is_indoor":         fmt.Sprintf("%v", place.IsIndoor),
			"with_child":        fmt.Sprintf("%v", place.WithChild),
			"with_pets":         fmt.Sprintf("%v", place.WithPets),
			"description":       place.Description,
		},
	}))
}

// editPlacePost — сохранение места
func (h *AdminHandler) editPlacePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin-error", h.aH(c, gin.H{"message": "Неверный ID"}))
		return
	}

	_, err = h.findPlaceByID(c, id)
	if err != nil {
		c.HTML(http.StatusNotFound, "admin-error", h.aH(c, gin.H{"message": err.Error()}))
		return
	}

	err = h.placeRepo.UpdatePlace(c.Request.Context(), id, map[string]any{
		"name":              c.PostForm("name"),
		"price":             c.PostForm("price"),
		"time":              func() float64 { var t float64; fmt.Sscanf(c.PostForm("time"), "%f", &t); return t }(),
		"types_of_movement": c.PostForm("types_of_movement"),
		"category":          c.PostForm("category"),
		"description":       c.PostForm("description"),
		"is_indoor":         c.PostForm("is_indoor") == "true",
		"with_child":        c.PostForm("with_child") == "true",
		"with_pets":         c.PostForm("with_pets") == "true",
	})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin-error", h.aH(c, gin.H{"message": err.Error()}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/table/places")
}

// deletePlace — удаление места
func (h *AdminHandler) deletePlace(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin-error", h.aH(c, gin.H{"message": "Неверный ID"}))
		return
	}

	err = h.placeRepo.DeletePlace(c.Request.Context(), id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin-error", h.aH(c, gin.H{"message": err.Error()}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/table/places")
}

// findPlaceByID — поиск места по ID
func (h *AdminHandler) findPlaceByID(c *gin.Context, id int) (*models.Place, error) {
	places, err := h.placeRepo.GetAll(c.Request.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to get places: %w", err)
	}

	for i := range places {
		if places[i].ID == id {
			return &places[i], nil
		}
	}

	return nil, fmt.Errorf("место не найдено")
}
