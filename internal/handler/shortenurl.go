package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyendong2003/bookmark-management/internal/service"
	"github.com/rs/zerolog/log"
)

type ShortenURL interface {
	ShortenURL(c *gin.Context)
	GetURL(c *gin.Context)
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
		log.Error().Err(err).Msg("Failed to shorten URL")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key})
}

// GetURL godoc
// @Summary      Redirect to original URL by code
// @Description  Retrieves the original URL associated with the given code. Browser clients are redirected (301); API/JSON clients receive the URL in the response body.
// @Tags         ShortenURL
// @Accept       json
// @Produce      json
// @Param        code   path      string  true  "Shortened URL code"
// @Success      301    {string}  string  "Redirects to the original URL (browser)"
// @Success      200    {object}  map[string]string  "Returns original URL (API/JSON client)"
// @Failure      400    {object}  map[string]string  "Code is required or not found"
// @Failure      500    {object}  map[string]string  "Internal server error"
// @Router       /link/redirect/{code} [get]
func (h *shortenURLHandler) GetURL(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code is required"})
		return
	}

	url, err := h.shortenURLService.GetURL(c, code)
	if err != nil {
		if errors.Is(err, service.ErrCodeNotExist) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Code not exist"})
			return
		}

		log.Error().Err(err).Msg("Failed to get URL")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Thêm dòng này để có response trả về khi test trên swagger UI hoặc API client, tránh lỗi CORS do fetch() tự động theo dõi redirects.
	// If client accepts JSON (e.g. Swagger UI, API clients), return URL in body
	// to avoid CORS issues caused by fetch() auto-following cross-origin redirects.
	// Browser navigation (Accept: text/html) gets the normal 301 redirect.
	if c.GetHeader("Accept") == "application/json" {
		c.JSON(http.StatusOK, gin.H{"url": url})
		return
	}

	// redirect browser client to url
	c.Redirect(http.StatusMovedPermanently, url)
}
