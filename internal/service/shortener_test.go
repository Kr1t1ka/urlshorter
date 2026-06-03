package service

import (
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/Kr1t1ka/shortUrl/internal/service/mocks"
)

func TestShorten(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStorage(ctrl)

	store.EXPECT().Get(gomock.Any()).Return("", false)
	store.EXPECT().Set(gomock.Any(), "https://example.com")

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
