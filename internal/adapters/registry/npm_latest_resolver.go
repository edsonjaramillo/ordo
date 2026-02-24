package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const defaultNPMPackageURL = "https://registry.npmjs.org"

type NPMLatestResolver struct {
	client   *http.Client
	endpoint string
}

func NewNPMLatestResolver() *NPMLatestResolver {
	return &NPMLatestResolver{
		client:   &http.Client{Timeout: defaultHTTPTimeout},
		endpoint: defaultNPMPackageURL,
	}
}

func (r *NPMLatestResolver) LatestVersion(ctx context.Context, packageName string) (string, error) {
	name := strings.TrimSpace(packageName)
	if name == "" {
		return "", fmt.Errorf("package name cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.endpoint+"/"+url.PathEscape(name), nil)
	if err != nil {
		return "", err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("npm package metadata status: %s", resp.Status)
	}

	var payload struct {
		DistTags map[string]string `json:"dist-tags"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}

	latest := strings.TrimSpace(payload.DistTags["latest"])
	if latest == "" {
		return "", fmt.Errorf("npm latest version not found for %s", name)
	}
	return latest, nil
}
