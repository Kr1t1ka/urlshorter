package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestEngine() *gin.Engine {
	r := gin.New()
	r.Use(GzipDecompress())
	r.Use(GzipCompress())
	return r
}

func gzipBody(t *testing.T, data string) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte(data))
	gz.Close()
	return &buf
}

func TestGzipDecompress(t *testing.T) {
	r := newTestEngine()
	r.POST("/", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.String(http.StatusOK, string(body))
	})

	buf := gzipBody(t, "hello")
	req := httptest.NewRequest(http.MethodPost, "/", buf)
	req.Header.Set("Content-Encoding", "gzip")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rr.Code, http.StatusOK)
	}
	if rr.Body.String() != "hello" {
		t.Errorf("got body %q, want %q", rr.Body.String(), "hello")
	}
}

func TestGzipDecompressInvalidBody(t *testing.T) {
	r := newTestEngine()
	r.POST("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("not gzip"))
	req.Header.Set("Content-Encoding", "gzip")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestGzipCompressJSON(t *testing.T) {
	r := newTestEngine()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"key": "value"})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rr.Code, http.StatusOK)
	}
	if rr.Header().Get("Content-Encoding") != "gzip" {
		t.Errorf("expected Content-Encoding: gzip, got %q", rr.Header().Get("Content-Encoding"))
	}

	gr, err := gzip.NewReader(rr.Body)
	if err != nil {
		t.Fatalf("response is not valid gzip: %v", err)
	}
	defer gr.Close()
	body, _ := io.ReadAll(gr)
	if !bytes.Contains(body, []byte("value")) {
		t.Errorf("decompressed body %q does not contain expected content", body)
	}
}

func TestGzipCompressNotAppliedForPlainText(t *testing.T) {
	r := newTestEngine()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "plain text")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Header().Get("Content-Encoding") == "gzip" {
		t.Error("expected no gzip compression for text/plain")
	}
	if rr.Body.String() != "plain text" {
		t.Errorf("got body %q, want %q", rr.Body.String(), "plain text")
	}
}

func TestGzipCompressNotAppliedWithoutHeader(t *testing.T) {
	r := newTestEngine()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"key": "value"})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Header().Get("Content-Encoding") == "gzip" {
		t.Error("expected no gzip when Accept-Encoding not set")
	}
}
