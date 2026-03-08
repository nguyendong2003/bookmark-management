package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyendong2003/bookmark-management/internal/service"
)

type Password interface {
	GenPass(c *gin.Context)
	GenPassForMux(w http.ResponseWriter, r *http.Request)
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
		c.String(http.StatusInternalServerError, "Failed to generate password")
		return
	}

	c.String(http.StatusOK, password)

	// c.JSON(http.StatusOK, gin.H{"password": password})
}

// Đoạn code này mô tả trường hợp nếu thay đổi gin bằng framework khác (ở đây là mux)
func (h *passwordHandler) GenPassForMux(w http.ResponseWriter, r *http.Request) {
	password, err := h.passwordService.GeneratePassword()
	if err != nil {
		http.Error(w, "Failed to generate password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(password))
}
