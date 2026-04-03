package handlers

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
	"github.com/ilushew/udmurtia-trip/backend/internal/services"
)

var errorTemplatePath = "auth-error"

type AuthHandler struct {
	userRepo *repository.UserRepository
	emailSvc *services.EmailService
}

func NewAuthHandler(userRepo *repository.UserRepository, emailSvc *services.EmailService) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		emailSvc: emailSvc,
	}
}

func (h *AuthHandler) ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register", nil)
}

func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	email := c.PostForm("email")

	if email == "" {
		c.HTML(http.StatusOK, errorTemplatePath, gin.H{"mesage": "Введите email"})
		return
	}
	user, err := h.userRepo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		user, err = h.userRepo.CreateUser(ctx, email, "")
		if err != nil {
			c.HTML(http.StatusOK, errorTemplatePath, gin.H{"message": "Ошибка регистрации"})
			return
		}
	}
	code := fmt.Sprintf("%06d", rand.IntN(1000000))
	expiresAt := time.Now().Add(15 * time.Minute)

	err = h.userRepo.SetVerificationCode(ctx, user.ID, code, expiresAt)
	if err != nil {
		c.HTML(http.StatusOK, errorTemplatePath, gin.H{"message": "Ошибка сохранения кода"})
		return
	}

	err = h.emailSvc.SendVerificationCode(email, code)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		c.HTML(http.StatusOK, errorTemplatePath, gin.H{"message": "Ошибка отправки письма"})
		return
	}
	c.HTML(http.StatusOK, "verify", gin.H{"email": email})
}

func (h *AuthHandler) VerifyCode(c *gin.Context) {
	ctx := c.Request.Context()
	email := c.PostForm("email")
	code := c.PostForm("code")

	user, err := h.userRepo.FindByEmail(ctx, email)
	if err != nil {
		c.HTML(http.StatusOK, errorTemplatePath, gin.H{"message": "Неверный email"})
		return
	}

	err = h.userRepo.VerifyCode(ctx, user.ID, code)
	if err != nil {
		c.HTML(http.StatusOK, errorTemplatePath, gin.H{"message": "Неверный или истёкший код"})
		return
	}

	err = h.userRepo.UpdateVerificationStatus(ctx, user.ID, true)
	if err != nil {
		c.HTML(http.StatusOK, errorTemplatePath, gin.H{"message": "Ошибка подтверждения"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID.String())
	session.Set("email", user.Email)
	err = session.Save()
	if err != nil {
		c.HTML(http.StatusOK, errorTemplatePath, gin.H{"message": "Ошибка создания сессии"})
		return
	}
	c.HTML(http.StatusOK, errorTemplatePath, gin.H{"message": "Ошибка - УСПЕХ"})
	// c.Header("HX-Redirect", "/")
	// c.AbortWithStatus(http.StatusOK)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Header("HX-Redirect", "/")
	c.Status(http.StatusNoContent)
}
