package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var ErrInvalidURL = errors.New("invalid url")

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

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Result string `json:"result"`
}

func (h *Handler) Register(r *gin.Engine) {
	r.POST("/", h.shortenHandler)
	r.POST("/api/shorten", h.apiShortenHandler)
	r.GET("/:id", h.redirectHandler)
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusBadRequest, "bad request")
	})
}

func validateURL(raw string) error {
	u, err := url.ParseRequestURI(raw)
	if err != nil || u.Host == "" {
		return fmt.Errorf("%w: %q", ErrInvalidURL, raw)
	}
	return nil
}

func (h *Handler) makeShortURL(originalURL string) (string, error) {
	if err := validateURL(originalURL); err != nil {
		return "", err
	}
	id, err := h.shortener.Shorten(originalURL)
	if err != nil {
		return "", fmt.Errorf("shorten: %w", err)
	}
	shortURL, err := url.JoinPath(h.baseURL, id)
	if err != nil {
		return "", fmt.Errorf("join path: %w", err)
	}
	return shortURL, nil
}

func (h *Handler) shortenHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	shortURL, err := h.makeShortURL(strings.TrimSpace(string(body)))
	if err != nil {
		if errors.Is(err, ErrInvalidURL) {
			c.String(http.StatusBadRequest, "invalid url")
			return
		}
		log.Printf("shorten error: %v", err)
		c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	c.String(http.StatusCreated, shortURL)
}

func (h *Handler) apiShortenHandler(c *gin.Context) {
	var req shortenRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	shortURL, err := h.makeShortURL(req.URL)
	if err != nil {
		if errors.Is(err, ErrInvalidURL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
			return
		}
		log.Printf("shorten error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, shortenResponse{Result: shortURL})
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
