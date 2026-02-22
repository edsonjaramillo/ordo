package execadapter

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"ordo/internal/domain"
)

func TestParseJSONDependenciesOutput(t *testing.T) {
	raw := `{"dependencies":{"typescript":{"version":"5.0.0"},"@types/node":{"version":"20.0.0"}}}`
	got, err := parseJSONDependenciesOutput(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != "@types/node" || got[1] != "typescript" {
		t.Fatalf("unexpected output: %#v", got)
	}
}

func TestParseYarnGlobalListOutput(t *testing.T) {
	raw := "{\"type\":\"tree\",\"data\":{\"type\":\"list\",\"trees\":[{\"name\":\"typescript@5.0.0\"},{\"name\":\"@types/node@20.0.0\"}]}}\n"
	got := parseYarnGlobalListOutput(raw)
	if len(got) != 2 || got[0] != "@types/node" || got[1] != "typescript" {
		t.Fatalf("unexpected output: %#v", got)
	}
}

func TestListGlobalCommand(t *testing.T) {
	tests := []struct {
		manager domain.PackageManager
		want0   string
		want1   string
	}{
		{manager: domain.ManagerNPM, want0: "npm", want1: "ls"},
		{manager: domain.ManagerPNPM, want0: "pnpm", want1: "ls"},
		{manager: domain.ManagerYarn, want0: "yarn", want1: "global"},
		{manager: domain.ManagerBun, want0: "bun", want1: "pm"},
	}

	for _, tc := range tests {
		got, err := listGlobalCommand(tc.manager)
		if err != nil {
			t.Fatalf("listGlobalCommand(%s) error = %v", tc.manager, err)
		}
		if len(got) < 2 || got[0] != tc.want0 || got[1] != tc.want1 {
			t.Fatalf("listGlobalCommand(%s) = %#v", tc.manager, got)
		}
	}
}

func TestListPackagesFromNodeModules(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "@types", "node"), 0o755); err != nil {
		t.Fatalf("mkdir scoped package: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "typescript"), 0o755); err != nil {
		t.Fatalf("mkdir package: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".bin"), 0o755); err != nil {
		t.Fatalf("mkdir dot dir: %v", err)
	}

	got, err := listPackagesFromNodeModules(root)
	if err != nil {
		t.Fatalf("listPackagesFromNodeModules() error = %v", err)
	}
	if len(got) != 2 || got[0] != "@types/node" || got[1] != "typescript" {
		t.Fatalf("unexpected output: %#v", got)
	}
}

func TestResolveGlobalStorePathsBunFromEnv(t *testing.T) {
	t.Setenv("BUN_INSTALL", "/tmp/custom-bun")
	r := NewRunner()

	paths, err := r.ResolveGlobalStorePaths(context.Background(), domain.ManagerBun)
	if err != nil {
		t.Fatalf("ResolveGlobalStorePaths() error = %v", err)
	}

	want := filepath.Clean("/tmp/custom-bun/install/global/node_modules")
	if len(paths) != 1 || paths[0] != want {
		t.Fatalf("paths = %#v, want [%q]", paths, want)
	}
}
