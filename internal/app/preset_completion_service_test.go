package app

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestPresetCompletionServicePresetNames(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	svc := NewPresetCompletionService(fakeConfigStore{
		content: []byte(`{
  "defaultPackageManager":"pnpm",
  "presets": {
    "eslint": {"devDependencies": ["eslint"]},
    "prettier": {"devDependencies": ["prettier"]}
  }
}`),
	})

	items, err := svc.PresetNames(context.Background(), "pre")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0] != "prettier" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestPresetCompletionServiceBuckets(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	svc := NewPresetCompletionService(fakeConfigStore{
		content: []byte(`{
  "defaultPackageManager":"pnpm",
  "presets": {
    "prettier": {
      "dependencies": [],
      "devDependencies": ["prettier"],
      "peerDependencies": ["   "],
      "optionalDependencies": ["prettier-plugin-tailwindcss"]
    }
  }
}`),
	})

	items, err := svc.Buckets(context.Background(), "prettier", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 || items[0] != "devDependencies" || items[1] != "optionalDependencies" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestPresetCompletionServiceBucketsUnknownPreset(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	svc := NewPresetCompletionService(fakeConfigStore{
		content: []byte(`{
  "defaultPackageManager":"pnpm",
  "presets": {
    "prettier": {"devDependencies": ["prettier"]}
  }
}`),
	})

	items, err := svc.Buckets(context.Background(), "unknown", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty items, got %#v", items)
	}
}

func TestPresetCompletionServiceBucketsPrefixFilter(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	svc := NewPresetCompletionService(fakeConfigStore{
		content: []byte(`{
  "defaultPackageManager":"pnpm",
  "presets": {
    "prettier": {
      "dependencies": ["prettier"],
      "devDependencies": ["prettier-plugin-tailwindcss"]
    }
  }
}`),
	})

	items, err := svc.Buckets(context.Background(), "prettier", "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0] != "devDependencies" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestPresetCompletionServiceBucketPackages(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	svc := NewPresetCompletionService(fakeConfigStore{
		content: []byte(`{
  "defaultPackageManager":"pnpm",
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

	items, err := svc.BucketPackages(context.Background(), "prettier", "devDependencies", "prettier-plugin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0] != "prettier-plugin-tailwindcss" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestPresetCompletionServiceMissingConfig(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	svc := NewPresetCompletionService(fakeConfigStore{err: os.ErrNotExist})

	_, err := svc.PresetNames(context.Background(), "")
	if !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("expected ErrConfigNotFound, got %v", err)
	}
}
