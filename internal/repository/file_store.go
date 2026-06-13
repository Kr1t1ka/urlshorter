package repository

import (
	"encoding/json"
	"fmt"
	"os"
)

type record struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileStore struct {
	*Store
	file    *os.File
	counter int
}

func NewFileStore(path string) (*FileStore, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open file store: %w", err)
	}

	fs := &FileStore{
		Store: NewStore(),
		file:  file,
	}

	decoder := json.NewDecoder(file)
	for decoder.More() {
		var r record
		if err := decoder.Decode(&r); err != nil {
			return nil, fmt.Errorf("decode record: %w", err)
		}
		fs.Store.Set(r.ShortURL, r.OriginalURL)
		fs.counter++
	}

	return fs, nil
}

func (fs *FileStore) Set(id, url string) error {
	if err := fs.Store.Set(id, url); err != nil {
		return err
	}

	fs.counter++
	r := record{
		UUID:        fmt.Sprintf("%d", fs.counter),
		ShortURL:    id,
		OriginalURL: url,
	}

	return json.NewEncoder(fs.file).Encode(r)
}
