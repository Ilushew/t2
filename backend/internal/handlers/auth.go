package handlers

import (
	"fmt"
	"log"
	"math/big"
	"crypto/rand"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
	"github.com/ilushew/udmurtia-trip/backend/internal/services"
)

var errorTemplatePath = "auth-error"

type AuthHandler struct {
	userRepo *repository.UserRepository
	emailSvc *services.EmailService
	codeSvc  *services.CodeService
}

func NewAuthHandler(userRepo *repository.UserRepository, emailSvc *services.EmailService, codeSvc *services.CodeService) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		emailSvc: emailSvc,
		codeSvc:  codeSvc,
	}
}

func generateCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
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
		user, err = h.userRepo.CreateUser(ctx, email)
		if err != nil {
			c.HTML(http.StatusOK, errorTemplatePath, gin.H{"message": "Ошибка регистрации"})
			return
		}
	}
	code, err := generateCode()
	if err != nil {
		log.Fatalf("Generate code error: %v", err)
	}

	err = h.codeSvc.SetCode(ctx, user.ID.String(), code)
	if err != nil {
		c.HTML(http.StatusOK, errorTemplatePath, gin.H{"message": "Ошибка сохранения кода"})
		return
	}
	go func() {
		err = h.emailSvc.SendVerificationCode(email, code)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	}()
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

	err = h.codeSvc.VerifyCode(ctx, user.ID.String(), code)
	if err != nil {
		c.HTML(http.StatusOK, errorTemplatePath, gin.H{"message": "Неверный или истёкший код"})
		return
	}

	err = h.userRepo.MarkVerified(ctx, user.ID)
	if err != nil {
		c.HTML(http.StatusOK, errorTemplatePath, gin.H{"message": "Ошибка обновления статуса"})
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
	c.Header("HX-Redirect", "/")
	c.AbortWithStatus(http.StatusOK)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
