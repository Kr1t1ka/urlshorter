package service

import "testing"

type mockStorage struct {
	data map[string]string
}

func newMockStorage() *mockStorage {
	return &mockStorage{data: map[string]string{}}
}

func (m *mockStorage) Set(id, url string) {
	m.data[id] = url
}

func (m *mockStorage) Get(id string) (string, bool) {
	url, ok := m.data[id]
	return url, ok
}

func TestShorten(t *testing.T) {
	store := newMockStorage()
	svc := NewShortener(store)

	id, err := svc.Shorten("https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty id")
	}

	stored, ok := store.Get(id)
	if !ok {
		t.Fatal("url not found in storage after Shorten")
	}
	if stored != "https://example.com" {
		t.Errorf("got %q, want %q", stored, "https://example.com")
	}
}

func TestResolve(t *testing.T) {
	store := newMockStorage()
	store.Set("abc123", "https://example.com")
	svc := NewShortener(store)

	got, ok := svc.Resolve("abc123")
	if !ok {
		t.Fatal("expected to find url")
	}
	if got != "https://example.com" {
		t.Errorf("got %q, want %q", got, "https://example.com")
	}
}

func TestResolveNotFound(t *testing.T) {
	svc := NewShortener(newMockStorage())

	_, ok := svc.Resolve("nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}
