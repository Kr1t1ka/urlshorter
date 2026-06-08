package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"

	"github.com/Kr1t1ka/shortUrl/internal/handler/mocks"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newEngine(svc shortenerService) *gin.Engine {
	r := gin.New()
	h := NewHandler(svc, "http://localhost:8080")
	h.Register(r)
	return r
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
		{"invalid url", "not-a-url", http.StatusBadRequest},
		{"url without host", "https:", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockshortenerService(ctrl)

			if tt.wantStatus == http.StatusCreated {
				svc.EXPECT().Shorten("https://example.com").Return("abc123", nil)
			}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()
			newEngine(svc).ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatus)
			}
		})
	}
}

func TestShortenHandlerServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	svc := mocks.NewMockshortenerService(ctrl)
	svc.EXPECT().Shorten("https://example.com").Return("", errors.New("storage error"))

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	rr := httptest.NewRecorder()
	newEngine(svc).ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusInternalServerError)
	}
}

func TestShortenHandlerResponseBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	svc := mocks.NewMockshortenerService(ctrl)
	svc.EXPECT().Shorten("https://example.com").Return("abc123", nil)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	rr := httptest.NewRecorder()
	newEngine(svc).ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("got status %d, want %d", rr.Code, http.StatusCreated)
	}
	if !strings.HasPrefix(rr.Body.String(), "http://localhost:8080/") {
		t.Errorf("unexpected response body: %q", rr.Body.String())
	}
}

func TestRedirectHandler(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		resolveURL   string
		resolveFound bool
		wantStatus   int
		wantLocation string
	}{
		{"existing id", "abc123", "https://example.com", true, http.StatusTemporaryRedirect, "https://example.com"},
		{"nonexistent id", "unknown", "", false, http.StatusBadRequest, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockshortenerService(ctrl)
			svc.EXPECT().Resolve(tt.id).Return(tt.resolveURL, tt.resolveFound)

			req := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)
			rr := httptest.NewRecorder()
			newEngine(svc).ServeHTTP(rr, req)

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
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockshortenerService(ctrl)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()
			newEngine(svc).ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("got status %d, want %d", rr.Code, http.StatusBadRequest)
			}
		})
	}
}
