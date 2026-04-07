package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
)

func RequireAdmin(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userIDStr := session.Get("user_id")
		if userIDStr == nil {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		user, err := userRepo.FindByID(c.Request.Context(), userID)
		if err != nil || !user.IsAdmin {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		// Передаём в шаблон автоматически
		c.Set("IsAuthenticated", true)
		c.Set("Email", user.Email)
		c.Next()
	}
}
