package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type IndexHandler struct {
}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (h *IndexHandler) ShowIndexPage(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	data := gin.H{
		"Title":            "Home",
		"IsAuthenticated":  userID != nil,
	}

	if userID != nil {
		data["Email"] = session.Get("email")
	}

	c.HTML(http.StatusOK, "index", data)
}