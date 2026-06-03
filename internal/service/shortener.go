package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

//go:generate mockgen -source=shortener.go -destination=mocks/mock_storage.go -package=mocks
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

const maxAttempts = 5

func (s *Shortener) Shorten(originalURL string) (string, error) {
	for range maxAttempts {
		id, err := generateID()
		if err != nil {
			return "", err
		}
		if _, exists := s.store.Get(id); !exists {
			s.store.Set(id, originalURL)
			return id, nil
		}
	}
	return "", errors.New("failed to generate unique id")
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
