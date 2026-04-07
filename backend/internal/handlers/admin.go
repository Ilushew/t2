package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
	"github.com/gin-contrib/sessions"
)

type AdminHandler struct {
	userRepo *repository.UserRepository
}

func NewAdminHandler(userRepo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{userRepo: userRepo}
}

// aH — хелпер для шаблонов: добавляет IsAuthenticated и Email
func (h *AdminHandler) aH(c *gin.Context, data gin.H) gin.H {
	session := sessions.Default(c)
	data["IsAuthenticated"] = true
	data["Email"] = session.Get("email")
	return data
}

// ShowTables — список доступных таблиц (только users)
func (h *AdminHandler) ShowTables(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-tables", h.aH(c, gin.H{
		"Title":  "Админ-панель",
		"Tables": []string{"users"},
	}))
}

// ViewTable — просмотр таблицы users
func (h *AdminHandler) ViewTable(c *gin.Context) {
	tableName := c.Param("table")
	if tableName != "users" {
		c.HTML(http.StatusNotFound, "admin-error", h.aH(c, gin.H{"message": "Таблица не найдена"}))
		return
	}

	users, err := h.userRepo.GetAll(c.Request.Context())
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin-error", h.aH(c, gin.H{"message": err.Error()}))
		return
	}

	columns := []struct{ Name, Type string }{
		{"id", "uuid"},
		{"email", "text"},
		{"is_verified", "bool"},
		{"is_admin", "bool"},
	}

	var rows [][]string
	for _, u := range users {
		rows = append(rows, []string{
			u.ID.String(),
			u.Email,
			fmt.Sprintf("%v", u.IsVerified),
			fmt.Sprintf("%v", u.IsAdmin),
		})
	}

	c.HTML(http.StatusOK, "admin-view", h.aH(c, gin.H{
		"Title":     "Таблица: users",
		"TableName": "users",
		"Columns":   columns,
		"Rows":      rows,
		"PKName":    "id",
	}))
}

// EditRowGet — форма редактирования пользователя
func (h *AdminHandler) EditRowGet(c *gin.Context) {
	tableName := c.Param("table")
	if tableName != "users" {
		c.HTML(http.StatusNotFound, "admin-error", h.aH(c, gin.H{"message": "Таблица не найдена"}))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin-error", h.aH(c, gin.H{"message": "Неверный ID"}))
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.HTML(http.StatusNotFound, "admin-error", h.aH(c, gin.H{"message": "Пользователь не найден"}))
		return
	}

	c.HTML(http.StatusOK, "admin-edit", h.aH(c, gin.H{
		"Title":     "Редактирование: users",
		"TableName": "users",
		"Columns":   []string{"id", "email", "is_verified", "is_admin"},
		"RowData": map[string]any{
			"id":          user.ID.String(),
			"email":       user.Email,
			"is_verified": fmt.Sprintf("%v", user.IsVerified),
			"is_admin":    fmt.Sprintf("%v", user.IsAdmin),
		},
	}))
}

// EditRowPost — сохранение изменений пользователя
func (h *AdminHandler) EditRowPost(c *gin.Context) {
	tableName := c.Param("table")
	if tableName != "users" {
		c.HTML(http.StatusNotFound, "admin-error", h.aH(c, gin.H{"message": "Таблица не найдена"}))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin-error", h.aH(c, gin.H{"message": "Неверный ID"}))
		return
	}

	err = h.userRepo.UpdateUser(c.Request.Context(), id, map[string]any{
		"email":       c.PostForm("email"),
		"is_verified": c.PostForm("is_verified") == "true",
		"is_admin":    c.PostForm("is_admin") == "true",
	})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin-error", h.aH(c, gin.H{"message": err.Error()}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/table/users")
}

// CreateUserGet — форма добавления пользователя
func (h *AdminHandler) CreateUserGet(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-create", h.aH(c, gin.H{
		"Title":     "Добавить пользователя",
		"TableName": "users",
	}))
}

// CreateUserPost — создание пользователя
func (h *AdminHandler) CreateUserPost(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.HTML(http.StatusBadRequest, "admin-error", h.aH(c, gin.H{"message": "Email обязателен"}))
		return
	}

	isVerified := c.PostForm("is_verified") == "on"
	isAdmin := c.PostForm("is_admin") == "on"

	_, err := h.userRepo.CreateUser(c.Request.Context(), email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin-error", h.aH(c, gin.H{"message": err.Error()}))
		return
	}

	// После создания — обновим флаги
	// Получаем пользователя по email
	user, err := h.userRepo.FindByEmail(c.Request.Context(), email)
	if err == nil && (isVerified || isAdmin) {
		updates := map[string]any{}
		if isVerified {
			updates["is_verified"] = true
		}
		if isAdmin {
			updates["is_admin"] = true
		}
		if len(updates) > 0 {
			h.userRepo.UpdateUser(c.Request.Context(), user.ID, updates)
		}
	}

	c.Redirect(http.StatusFound, "/admin/table/users")
}

// DeleteRow — удаление пользователя
func (h *AdminHandler) DeleteRow(c *gin.Context) {
	tableName := c.Param("table")
	if tableName != "users" {
		c.HTML(http.StatusNotFound, "admin-error", h.aH(c, gin.H{"message": "Таблица не найдена"}))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin-error", h.aH(c, gin.H{"message": "Неверный ID"}))
		return
	}

	err = h.userRepo.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin-error", h.aH(c, gin.H{"message": err.Error()}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/table/users")
}
