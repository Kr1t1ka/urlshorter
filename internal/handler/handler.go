package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/Kr1t1ka/shortUrl/internal/service"
)

type Handler struct {
	shortener *service.Shortener
}

func NewHandler(shortener *service.Shortener) *Handler {
	return &Handler{shortener: shortener}
}

func (h *Handler) Route(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/":
		h.shortenHandler(w, r)
	case r.Method == http.MethodGet && r.URL.Path != "/":
		h.redirectHandler(w, r)
	default:
		http.Error(w, "bad request", http.StatusBadRequest)
	}
}

func (h *Handler) shortenHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(string(body))
	if len(originalURL) == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	id, err := h.shortener.Shorten(originalURL)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("http://localhost:8080/" + id))
}

func (h *Handler) redirectHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	original, ok := h.shortener.Resolve(id)
	if !ok {
		http.Error(w, "not found", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, original, http.StatusTemporaryRedirect)
}
