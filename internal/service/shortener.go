package service

import (
	"crypto/rand"
	"encoding/base64"
)

type Storage interface {
	Set(id, url string)
	Get(id string) (string, bool)
}

type Shortener struct {
	store Storage
}

func NewShortener(store Storage) *Shortener {
	return &Shortener{store: store}
}

func (s *Shortener) Shorten(originalURL string) (string, error) {
	id, err := generateID()
	if err != nil {
		return "", err
	}
	s.store.Set(id, originalURL)
	return id, nil
}

func (s *Shortener) Resolve(id string) (string, bool) {
	return s.store.Get(id)
}

func generateID() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
