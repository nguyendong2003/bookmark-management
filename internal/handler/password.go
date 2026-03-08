package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyendong2003/bookmark-management/internal/service"
)

type Password interface {
	GenPass(c *gin.Context)
}

type passwordHandler struct {
	passwordService service.Password
}

func NewPassword(passwordService service.Password) Password {
	return &passwordHandler{
		passwordService: passwordService,
	}
}

func (h *passwordHandler) GenPass(c *gin.Context) {
	password, err := h.passwordService.GeneratePassword()
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate password"})
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, password)

	// c.JSON(http.StatusOK, gin.H{"password": password})
}
