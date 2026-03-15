package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyendong2003/bookmark-management/internal/service"
)

type ShortenURL interface {
	ShortenURL(c *gin.Context)
}

type shortenURLHandler struct {
	shortenURLService service.ShortenURL
}

type shortenURLRequest struct {
	URL string `json:"url" binding:"required,url"`
}

func NewShortenURL(shortenURLService service.ShortenURL) ShortenURL {
	return &shortenURLHandler{
		shortenURLService: shortenURLService,
	}
}

// ShortenURL godoc
// @Summary Shorten a URL
// @Description Create a short key for the provided URL
// @Tags ShortenURL
// @Accept json
// @Produce json
// @Param request body shortenURLRequest true "URL to shorten"
// @Success 200 {object} map[string]string "Shortened key"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /link/shorten [post]
func (h *shortenURLHandler) ShortenURL(c *gin.Context) {
	req := &shortenURLRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	key, err := h.shortenURLService.ShortenURL(c, req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key})
}
