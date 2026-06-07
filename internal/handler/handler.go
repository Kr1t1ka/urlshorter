package handler

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock_shortener_service.go -package=mocks
type shortenerService interface {
	Shorten(url string) (string, error)
	Resolve(id string) (string, bool)
}

type Handler struct {
	shortener shortenerService
	baseURL   string
}

func NewHandler(shortener shortenerService, baseURL string) *Handler {
	return &Handler{shortener: shortener, baseURL: baseURL}
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

	if u, err := url.ParseRequestURI(originalURL); err != nil || u.Host == "" {
		c.String(http.StatusBadRequest, "invalid url")
		return
	}

	id, err := h.shortener.Shorten(originalURL)
	if err != nil {
		log.Printf("shorten error: %v", err)
		c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	shortURL, err := url.JoinPath(h.baseURL, id)
	if err != nil {
		log.Printf("url join error: %v", err)
		c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	c.String(http.StatusCreated, shortURL)
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
