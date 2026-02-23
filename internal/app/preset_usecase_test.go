package app

import (
	"context"
	"errors"
	"os"
	"testing"
)

type fakeConfigStore struct {
	content []byte
	err     error
}

func (f fakeConfigStore) MkdirAll(string, os.FileMode) error { return nil }
func (f fakeConfigStore) Exists(string) (bool, error)        { return false, nil }
func (f fakeConfigStore) WriteFile(string, []byte, os.FileMode) error {
	return nil
}

func (f fakeConfigStore) ReadFile(string) ([]byte, error) {
	if f.err != nil {
		return nil, f.err
	}
	return append([]byte(nil), f.content...), nil
}

func TestPresetUseCaseBucketInstall(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewPresetUseCase(discovery, runner, fakeConfigStore{
		content: []byte(`{
  "defaultPackageManager": "pnpm",
  "presets": {
    "prettier": {
      "devDependencies": [
        "prettier",
        "@ianvs/prettier-plugin-sort-imports"
      ]
    }
  }
}`),
	})

	err := uc.Run(context.Background(), PresetRequest{
		Preset: "prettier",
		Bucket: "devDependencies",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if runner.dir != "." {
		t.Fatalf("expected root dir '.', got %s", runner.dir)
	}
	want := []string{"pnpm", "add", "--save-dev", "prettier", "@ianvs/prettier-plugin-sort-imports"}
	if len(runner.argv) != len(want) {
		t.Fatalf("unexpected argv len: %#v", runner.argv)
	}
	for i := range want {
		if runner.argv[i] != want[i] {
			t.Fatalf("argv[%d] = %q, want %q", i, runner.argv[i], want[i])
		}
	}
}

func TestPresetUseCaseWorkspaceWithFilter(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewPresetUseCase(discovery, runner, fakeConfigStore{
		content: []byte(`{
  "defaultPackageManager": "pnpm",
  "presets": {
    "prettier": {
      "devDependencies": [
        "prettier",
        "prettier-plugin-tailwindcss"
      ]
    }
  }
}`),
	})

	err := uc.Run(context.Background(), PresetRequest{
		Preset:    "prettier",
		Bucket:    "devDependencies",
		Packages:  []string{"prettier-plugin-tailwindcss"},
		Workspace: "ui",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if runner.dir != "packages/ui" {
		t.Fatalf("expected workspace dir packages/ui, got %s", runner.dir)
	}
	want := []string{"pnpm", "add", "--save-dev", "prettier-plugin-tailwindcss"}
	if len(runner.argv) != len(want) {
		t.Fatalf("unexpected argv len: %#v", runner.argv)
	}
	for i := range want {
		if runner.argv[i] != want[i] {
			t.Fatalf("argv[%d] = %q, want %q", i, runner.argv[i], want[i])
		}
	}
}

func TestPresetUseCaseUnknownFilteredPackage(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewPresetUseCase(discovery, runner, fakeConfigStore{
		content: []byte(`{
  "defaultPackageManager": "pnpm",
  "presets": {
    "prettier": {
      "devDependencies": ["prettier"]
    }
  }
}`),
	})

	err := uc.Run(context.Background(), PresetRequest{
		Preset:   "prettier",
		Bucket:   "devDependencies",
		Packages: []string{"prettier-plugin-tailwindcss"},
	})
	if !errors.Is(err, ErrPresetPackageNotFound) {
		t.Fatalf("expected ErrPresetPackageNotFound, got %v", err)
	}
}

func TestPresetUseCaseMissingConfig(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewPresetUseCase(discovery, runner, fakeConfigStore{err: os.ErrNotExist})

	err := uc.Run(context.Background(), PresetRequest{
		Preset: "prettier",
		Bucket: "devDependencies",
	})
	if !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("expected ErrConfigNotFound, got %v", err)
	}
}

func TestPresetUseCaseMissingPreset(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewPresetUseCase(discovery, runner, fakeConfigStore{
		content: []byte(`{"defaultPackageManager":"pnpm","presets":{}}`),
	})

	err := uc.Run(context.Background(), PresetRequest{
		Preset: "prettier",
		Bucket: "devDependencies",
	})
	if !errors.Is(err, ErrPresetNotFound) {
		t.Fatalf("expected ErrPresetNotFound, got %v", err)
	}
}

func TestPresetUseCaseMissingBucket(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewPresetUseCase(discovery, runner, fakeConfigStore{
		content: []byte(`{
  "defaultPackageManager": "pnpm",
  "presets": {"prettier": {}}
}`),
	})

	err := uc.Run(context.Background(), PresetRequest{
		Preset: "prettier",
		Bucket: "devDependencies",
	})
	if !errors.Is(err, ErrPresetBucketNotFound) {
		t.Fatalf("expected ErrPresetBucketNotFound, got %v", err)
	}
}
