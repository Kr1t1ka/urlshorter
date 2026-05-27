package repository

type Store struct {
	data map[string]string
}

func NewStore() *Store {
	return &Store{data: map[string]string{}}
}

func (s *Store) Set(id, url string) {
	s.data[id] = url
}

func (s *Store) Get(id string) (string, bool) {
	url, ok := s.data[id]
	return url, ok
}
