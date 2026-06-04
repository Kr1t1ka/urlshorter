package repository

import (
	"errors"
	"sync"
)

type Store struct {
	mu   sync.Mutex
	data map[string]string
}

func NewStore() *Store {
	return &Store{data: map[string]string{}}
}

func (s *Store) Set(id, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[id]; exists {
		return errors.New("already exists")
	}
	s.data[id] = url
	return nil
}

func (s *Store) Get(id string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	url, ok := s.data[id]
	return url, ok
}
