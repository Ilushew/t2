package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
)

type ProfileHandler struct {
	userRepo *repository.UserRepository
}

func NewProfileHandler(userRepo *repository.UserRepository) *ProfileHandler {
	return &ProfileHandler{
		userRepo: userRepo,
	}
}

func (h *ProfileHandler) ShowProfilePage(c *gin.Context) {
	session := sessions.Default(c)
	userIDStr := session.Get("user_id")

	if userIDStr == nil {
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.HTML(http.StatusOK, "auth-error", gin.H{"message": "Пользователь не найден"})
		return
	}

	data := gin.H{
		"Title":           "Профиль",
		"Email":           user.Email,
		"IsVerified":      user.IsVerified,
		"IsAuthenticated": true,
	}

	c.HTML(http.StatusOK, "profile", data)
}
