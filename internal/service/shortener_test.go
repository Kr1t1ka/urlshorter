package service

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/Kr1t1ka/shortUrl/internal/repository"
	"github.com/Kr1t1ka/shortUrl/internal/service/mocks"
)

func TestShorten(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStorage(ctrl)

	store.EXPECT().Set(gomock.Any(), "https://example.com").Return(nil)

	svc := NewShortener(store)

	id, err := svc.Shorten("https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty id")
	}
}

func TestResolve(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStorage(ctrl)
	store.EXPECT().Get("abc123").Return("https://example.com", true)

	svc := NewShortener(store)

	got, ok := svc.Resolve("abc123")
	if !ok {
		t.Fatal("expected to find url")
	}
	if got != "https://example.com" {
		t.Errorf("got %q, want %q", got, "https://example.com")
	}
}

func TestShortenStorageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStorage(ctrl)
	store.EXPECT().Set(gomock.Any(), "https://example.com").Return(errors.New("storage error"))

	svc := NewShortener(store)

	_, err := svc.Shorten("https://example.com")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestShortenIDCollision(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStorage(ctrl)
	store.EXPECT().Set(gomock.Any(), "https://example.com").Return(repository.ErrIDExists).Times(maxAttempts)

	svc := NewShortener(store)

	_, err := svc.Shorten("https://example.com")
	if err == nil {
		t.Fatal("expected error after max attempts")
	}
}

func TestResolveNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStorage(ctrl)
	store.EXPECT().Get("nonexistent").Return("", false)

	svc := NewShortener(store)

	_, ok := svc.Resolve("nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}
