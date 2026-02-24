package catalog

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	fsadapter "ordo/internal/adapters/fs"
)

func TestManifestStoreRewriteCatalogReferencesExistingOnly(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "packages/ui/package.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	content := []byte(`{
  "name": "@repo/ui",
  "scripts": { "build": "tsup" },
  "dependencies": { "react": "^19.0.0" },
  "devDependencies": { "typescript": "^5.0.0" }
}
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store := NewManifestStore(root, fsadapter.NewConfigStore())
	err := store.RewriteCatalogReferencesExistingOnly(context.Background(), "packages/ui", "", []string{"react", "zod"})
	if err != nil {
		t.Fatalf("RewriteCatalogReferencesExistingOnly() error = %v", err)
	}

	manifest := readManifest(t, path)
	deps := asStringMap(t, manifest["dependencies"])
	if deps["react"] != "catalog:" {
		t.Fatalf("dependencies.react = %q, want catalog:", deps["react"])
	}
	if _, ok := deps["zod"]; ok {
		t.Fatalf("dependencies should not add zod, got %#v", deps)
	}
	if _, ok := manifest["scripts"]; !ok {
		t.Fatalf("scripts field missing after rewrite: %#v", manifest)
	}
}

func TestManifestStoreRewriteCatalogReferencesAddsMissing(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "packages/ui/package.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	content := []byte(`{
  "name": "@repo/ui",
  "dependencies": {}
}
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store := NewManifestStore(root, fsadapter.NewConfigStore())
	err := store.RewriteCatalogReferences(context.Background(), "packages/ui", "", []string{"react"})
	if err != nil {
		t.Fatalf("RewriteCatalogReferences() error = %v", err)
	}

	manifest := readManifest(t, path)
	deps := asStringMap(t, manifest["dependencies"])
	if deps["react"] != "catalog:" {
		t.Fatalf("dependencies.react = %q, want catalog:", deps["react"])
	}
}

func readManifest(t *testing.T, path string) map[string]any {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	manifest := map[string]any{}
	if err := json.Unmarshal(content, &manifest); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	return manifest
}

func asStringMap(t *testing.T, raw any) map[string]string {
	t.Helper()
	out := map[string]string{}
	items, ok := raw.(map[string]any)
	if !ok {
		return out
	}
	for key, value := range items {
		str, ok := value.(string)
		if !ok {
			t.Fatalf("value for %q is not string: %#v", key, value)
		}
		out[key] = str
	}
	return out
}
