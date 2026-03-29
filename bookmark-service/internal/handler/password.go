package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyendong2003/bookmark-management/internal/service"
	"github.com/rs/zerolog/log"
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

// GenPass godoc
// @Summary Generate a random password
// @Description Generate a random password with specified criteria
// @Tags Password
// @Produce plain
// @Success 200 {string} string "Generated password"
// @Failure 500 {string} string "Failed to generate password"
// @Router /gen-pass [get]
func (h *passwordHandler) GenPass(c *gin.Context) {
	password, err := h.passwordService.GeneratePassword()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate password")

		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate password"})
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, password)

	// c.JSON(http.StatusOK, gin.H{"password": password})
}
