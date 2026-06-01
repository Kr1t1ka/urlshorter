package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Kr1t1ka/shortUrl/internal/config"
	"github.com/Kr1t1ka/shortUrl/internal/handler"
	"github.com/Kr1t1ka/shortUrl/internal/repository"
	"github.com/Kr1t1ka/shortUrl/internal/service"
)

func main() {
	cfg := config.New()

	store := repository.NewStore()
	shortener := service.NewShortener(store)
	h := handler.NewHandler(shortener, cfg.BaseURL)

	r := gin.Default()
	h.Register(r)

	log.Fatal(r.Run(cfg.ServerAddr))
}
