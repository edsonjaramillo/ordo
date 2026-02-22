package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultNPMSearchURL = "https://registry.npmjs.org/-/v1/search"
	defaultCacheTTL     = time.Minute
	defaultHTTPTimeout  = 1500 * time.Millisecond
)

type NPMSuggestor struct {
	client   *http.Client
	endpoint string
	ttl      time.Duration

	mu    sync.Mutex
	cache map[string]cacheEntry
}

type cacheEntry struct {
	expiresAt time.Time
	items     []string
}

func NewNPMSuggestor() *NPMSuggestor {
	return &NPMSuggestor{
		client:   &http.Client{Timeout: defaultHTTPTimeout},
		endpoint: defaultNPMSearchURL,
		ttl:      defaultCacheTTL,
		cache:    map[string]cacheEntry{},
	}
}

func (s *NPMSuggestor) Suggest(ctx context.Context, prefix string, limit int) ([]string, error) {
	p := strings.TrimSpace(prefix)
	if p == "" {
		return nil, nil
	}
	if limit <= 0 {
		limit = 20
	}

	if items, ok := s.cached(p); ok {
		return items, nil
	}

	values := url.Values{}
	values.Set("text", p)
	values.Set("size", strconv.Itoa(limit))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.endpoint+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("npm search status: %s", resp.Status)
	}

	var payload struct {
		Objects []struct {
			Package struct {
				Name string `json:"name"`
			} `json:"package"`
		} `json:"objects"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	items := make([]string, 0, len(payload.Objects))
	for _, obj := range payload.Objects {
		name := strings.TrimSpace(obj.Package.Name)
		if name == "" {
			continue
		}
		if !strings.HasPrefix(name, p) {
			continue
		}
		items = append(items, name)
	}
	sort.Strings(items)
	items = unique(items)
	s.store(p, items)
	return items, nil
}

func unique(items []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func (s *NPMSuggestor) cached(prefix string) ([]string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.cache[prefix]
	if !ok {
		return nil, false
	}
	if time.Now().After(entry.expiresAt) {
		delete(s.cache, prefix)
		return nil, false
	}
	return append([]string(nil), entry.items...), true
}

func (s *NPMSuggestor) store(prefix string, items []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[prefix] = cacheEntry{
		expiresAt: time.Now().Add(s.ttl),
		items:     append([]string(nil), items...),
	}
}
