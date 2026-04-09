package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// viewUsers — просмотр таблицы users
func (h *AdminHandler) viewUsers(c *gin.Context) {
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

// editUserGet — форма редактирования пользователя
func (h *AdminHandler) editUserGet(c *gin.Context) {
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

// editUserPost — сохранение изменений пользователя
func (h *AdminHandler) editUserPost(c *gin.Context) {
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

// createUserGet — форма добавления пользователя
func (h *AdminHandler) createUserGet(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-create", h.aH(c, gin.H{
		"Title":     "Добавить пользователя",
		"TableName": "users",
	}))
}

// createUserPost — создание пользователя
func (h *AdminHandler) createUserPost(c *gin.Context) {
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

// deleteUser — удаление пользователя
func (h *AdminHandler) deleteUser(c *gin.Context) {
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
