package main

import (
	"io"
	"log"

	"github.com/Kr1t1ka/shortUrl/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Kr1t1ka/shortUrl/internal/config"
	"github.com/Kr1t1ka/shortUrl/internal/handler"
	"github.com/Kr1t1ka/shortUrl/internal/repository"
	"github.com/Kr1t1ka/shortUrl/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer func() { _ = logger.Sync() }()

	var store service.Storage
	if cfg.FileStoragePath != "" {
		fs, err := repository.NewFileStore(cfg.FileStoragePath)
		if err != nil {
			return err
		}
		store = fs
	} else {
		store = repository.NewStore()
	}
	if c, ok := store.(io.Closer); ok {
		defer func() { _ = c.Close() }()
	}

	shortener := service.NewShortener(store)
	h := handler.NewHandler(shortener, cfg.BaseURL)

	r := gin.New()
	r.Use(middleware.Logger(logger))
	r.Use(middleware.GzipDecompress())
	r.Use(middleware.GzipCompress())
	h.Register(r)

	return r.Run(cfg.ServerAddr)
}
