package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Kr1t1ka/shortUrl/internal/service"
)

type mockStorage struct {
	data map[string]string
}

func newMockStorage() *mockStorage {
	return &mockStorage{data: map[string]string{}}
}

func (m *mockStorage) Set(id, url string) {
	m.data[id] = url
}

func (m *mockStorage) Get(id string) (string, bool) {
	url, ok := m.data[id]
	return url, ok
}

func newTestHandler() *Handler {
	return NewHandler(service.NewShortener(newMockStorage()))
}

func TestShortenHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"valid url", "https://example.com", http.StatusCreated},
		{"empty body", "", http.StatusBadRequest},
		{"whitespace only", "   ", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()

			newTestHandler().Route(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatus)
			}
		})
	}
}

func TestShortenHandlerResponseBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	rr := httptest.NewRecorder()

	newTestHandler().Route(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("got status %d, want %d", rr.Code, http.StatusCreated)
	}
	if !strings.HasPrefix(rr.Body.String(), "http://localhost:8080/") {
		t.Errorf("unexpected response body: %q", rr.Body.String())
	}
}

func TestRedirectHandler(t *testing.T) {
	store := newMockStorage()
	store.Set("abc123", "https://example.com")
	h := NewHandler(service.NewShortener(store))

	tests := []struct {
		name         string
		id           string
		wantStatus   int
		wantLocation string
	}{
		{"existing id", "abc123", http.StatusTemporaryRedirect, "https://example.com"},
		{"nonexistent id", "unknown", http.StatusBadRequest, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)
			rr := httptest.NewRecorder()

			h.Route(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatus)
			}
			if tt.wantLocation != "" {
				if loc := rr.Header().Get("Location"); loc != tt.wantLocation {
					t.Errorf("got Location %q, want %q", loc, tt.wantLocation)
				}
			}
		})
	}
}

func TestRouteInvalidRequests(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"GET root", http.MethodGet, "/"},
		{"DELETE root", http.MethodDelete, "/"},
		{"PUT root", http.MethodPut, "/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			newTestHandler().Route(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("got status %d, want %d", rr.Code, http.StatusBadRequest)
			}
		})
	}
}
