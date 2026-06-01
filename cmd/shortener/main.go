package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Kr1t1ka/shortUrl/internal/handler"
	"github.com/Kr1t1ka/shortUrl/internal/repository"
	"github.com/Kr1t1ka/shortUrl/internal/service"
)

func main() {
	store := repository.NewStore()
	shortener := service.NewShortener(store)
	h := handler.NewHandler(shortener)

	r := gin.Default()
	h.Register(r)

	log.Fatal(r.Run(":8080"))
}
