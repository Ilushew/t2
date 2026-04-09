package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
)

type AdminHandler struct {
	userRepo  *repository.UserRepository
	placeRepo *repository.PlaceRepository
}

func NewAdminHandler(userRepo *repository.UserRepository, placeRepo *repository.PlaceRepository) *AdminHandler {
	return &AdminHandler{
		userRepo:  userRepo,
		placeRepo: placeRepo,
	}
}

// aH — хелпер для шаблонов: добавляет IsAuthenticated и Email
func (h *AdminHandler) aH(c *gin.Context, data gin.H) gin.H {
	session := sessions.Default(c)
	data["IsAuthenticated"] = true
	data["Email"] = session.Get("email")
	return data
}

// ShowTables — список доступных таблиц
func (h *AdminHandler) ShowTables(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-tables", h.aH(c, gin.H{
		"Title":  "Админ-панель",
		"Tables": []string{"users", "places"},
	}))
}

// ViewTable — маршрутизатор просмотра таблицы
func (h *AdminHandler) ViewTable(c *gin.Context) {
	tableName := c.Param("table")

	switch tableName {
	case "users":
		h.viewUsers(c)
	case "places":
		h.viewPlaces(c)
	default:
		c.HTML(http.StatusNotFound, "admin-error", h.aH(c, gin.H{"message": "Таблица не найдена"}))
	}
}

// EditRowGet — маршрутизатор формы редактирования
func (h *AdminHandler) EditRowGet(c *gin.Context) {
	tableName := c.Param("table")

	switch tableName {
	case "users":
		h.editUserGet(c)
	case "places":
		h.editPlaceGet(c)
	default:
		c.HTML(http.StatusNotFound, "admin-error", h.aH(c, gin.H{"message": "Таблица не найдена"}))
	}
}

// EditRowPost — маршрутизатор сохранения изменений
func (h *AdminHandler) EditRowPost(c *gin.Context) {
	tableName := c.Param("table")

	switch tableName {
	case "users":
		h.editUserPost(c)
	case "places":
		h.editPlacePost(c)
	default:
		c.HTML(http.StatusNotFound, "admin-error", h.aH(c, gin.H{"message": "Таблица не найдена"}))
	}
}

// CreateUserGet — маршрутизатор формы создания
func (h *AdminHandler) CreateUserGet(c *gin.Context) {
	tableName := c.Param("table")

	switch tableName {
	case "users":
		h.createUserGet(c)
	case "places":
		h.createPlaceGet(c)
	default:
		c.HTML(http.StatusNotFound, "admin-error", h.aH(c, gin.H{"message": "Таблица не найдена"}))
	}
}

// CreateUserPost — маршрутизатор создания записей
func (h *AdminHandler) CreateUserPost(c *gin.Context) {
	tableName := c.Param("table")

	switch tableName {
	case "users":
		h.createUserPost(c)
	case "places":
		h.createPlacePost(c)
	default:
		c.HTML(http.StatusNotFound, "admin-error", h.aH(c, gin.H{"message": "Таблица не найдена"}))
	}
}

// DeleteRow — маршрутизатор удаления записей
func (h *AdminHandler) DeleteRow(c *gin.Context) {
	tableName := c.Param("table")

	switch tableName {
	case "users":
		h.deleteUser(c)
	case "places":
		h.deletePlace(c)
	default:
		c.HTML(http.StatusNotFound, "admin-error", h.aH(c, gin.H{"message": "Таблица не найдена"}))
	}
}
