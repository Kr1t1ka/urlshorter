package main

import (
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
	cfg := config.New()
	logger, _ := zap.NewProduction()

	store := repository.NewStore()
	shortener := service.NewShortener(store)
	h := handler.NewHandler(shortener, cfg.BaseURL)

	r := gin.New()
	r.Use(middleware.Logger(logger))
	r.Use(middleware.GzipDecompress())
	r.Use(middleware.GzipCompress())
	h.Register(r)

	log.Fatal(r.Run(cfg.ServerAddr))
}
