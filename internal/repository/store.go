package repository

import "sync"

type Store struct {
	mu   sync.Mutex
	data map[string]string
}

func NewStore() *Store {
	return &Store{data: map[string]string{}}
}

func (s *Store) Set(id, url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = url
}

func (s *Store) Get(id string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	url, ok := s.data[id]
	return url, ok
}
