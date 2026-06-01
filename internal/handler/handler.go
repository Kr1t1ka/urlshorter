package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Kr1t1ka/shortUrl/internal/service"
)

type Handler struct {
	shortener *service.Shortener
}

func NewHandler(shortener *service.Shortener) *Handler {
	return &Handler{shortener: shortener}
}

func (h *Handler) Register(r *gin.Engine) {
	r.POST("/", h.shortenHandler)
	r.GET("/:id", h.redirectHandler)
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusBadRequest, "bad request")
	})
}

func (h *Handler) shortenHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	originalURL := strings.TrimSpace(string(body))
	if len(originalURL) == 0 {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	id, err := h.shortener.Shorten(originalURL)
	if err != nil {
		c.String(http.StatusInternalServerError, "internal error")
		return
	}

	c.String(http.StatusCreated, "http://localhost:8080/"+id)
}

func (h *Handler) redirectHandler(c *gin.Context) {
	id := c.Param("id")

	original, ok := h.shortener.Resolve(id)
	if !ok {
		c.String(http.StatusBadRequest, "not found")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, original)
}
