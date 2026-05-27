package main

import (
	"log"
	"net/http"

	"github.com/Kr1t1ka/shortUrl/internal/handler"
	"github.com/Kr1t1ka/shortUrl/internal/repository"
	"github.com/Kr1t1ka/shortUrl/internal/service"
)

func main() {
	store := repository.NewStore()
	shortener := service.NewShortener(store)
	h := handler.NewHandler(shortener)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Route)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
