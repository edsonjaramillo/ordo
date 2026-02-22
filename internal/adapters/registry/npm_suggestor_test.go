package registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNPMSuggestorSuggest(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		_, _ = w.Write([]byte(`{"objects":[{"package":{"name":"react"}},{"package":{"name":"react-dom"}},{"package":{"name":"vite"}}]}`))
	}))
	defer srv.Close()

	s := &NPMSuggestor{
		client:   srv.Client(),
		endpoint: srv.URL,
		ttl:      time.Minute,
		cache:    map[string]cacheEntry{},
	}

	items, err := s.Suggest(context.Background(), "react", 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 || items[0] != "react" || items[1] != "react-dom" {
		t.Fatalf("unexpected items: %#v", items)
	}
	if hits != 1 {
		t.Fatalf("expected 1 hit, got %d", hits)
	}

	again, err := s.Suggest(context.Background(), "react", 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(again) != 2 {
		t.Fatalf("unexpected cached items: %#v", again)
	}
	if hits != 1 {
		t.Fatalf("expected cache hit with 1 server call, got %d", hits)
	}
}
