package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GzipDecompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Encoding") == "gzip" {
			gr, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			defer gr.Close()
			c.Request.Body = io.NopCloser(gr)
			c.Request.Header.Del("Content-Encoding")
		}
		c.Next()
	}
}

func GzipCompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			c.Next()
			return
		}

		gw := &gzipResponseWriter{ResponseWriter: c.Writer, gz: gz}
		c.Writer = gw
		defer func() {
			if gw.used {
				_ = gz.Close()
			}
		}()
		c.Next()
	}
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	gz        *gzip.Writer
	used      bool
	committed bool
}

func (w *gzipResponseWriter) decide() {
	if w.committed {
		return
	}
	ct := w.Header().Get("Content-Type")
	if strings.Contains(ct, "application/json") || strings.Contains(ct, "text/html") {
		w.used = true
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")
		w.committed = true
	}
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	w.decide()
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	w.decide()
	if w.used {
		return w.gz.Write(data)
	}
	return w.ResponseWriter.Write(data)
}

func (w *gzipResponseWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}
