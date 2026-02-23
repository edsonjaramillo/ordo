package app

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fsadapter "ordo/internal/adapters/fs"
)

func TestInitUseCaseWritesConfigInXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	uc := NewInitUseCase(fsadapter.NewConfigStore())
	err := uc.Run(context.Background(), InitRequest{DefaultPackageManager: "pnpm"})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	configPath, err := filepath.Abs(filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "ordo", "ordo.json"))
	if err != nil {
		t.Fatalf("Abs() error = %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if got := string(content); got != "{\n  \"defaultPackageManager\": \"pnpm\"\n}\n" {
		t.Fatalf("config content = %q, want JSON with defaultPackageManager", got)
	}
}

func TestInitUseCaseFallsBackToHomeDotConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")

	uc := NewInitUseCase(fsadapter.NewConfigStore())
	err := uc.Run(context.Background(), InitRequest{DefaultPackageManager: "npm"})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	configPath := filepath.Join(home, ".config", "ordo", "ordo.json")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("Stat(%s) error = %v", configPath, err)
	}
}

func TestInitUseCaseFailsWhenConfigAlreadyExists(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)

	configDir := filepath.Join(xdg, "ordo")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	configPath := filepath.Join(configDir, "ordo.json")
	if err := os.WriteFile(configPath, []byte("{}\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	uc := NewInitUseCase(fsadapter.NewConfigStore())
	err := uc.Run(context.Background(), InitRequest{DefaultPackageManager: "yarn"})
	if !errors.Is(err, ErrConfigAlreadyExists) {
		t.Fatalf("Run() error = %v, want ErrConfigAlreadyExists", err)
	}
}

func TestInitUseCaseRejectsInvalidDefaultPackageManager(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	uc := NewInitUseCase(fsadapter.NewConfigStore())
	err := uc.Run(context.Background(), InitRequest{DefaultPackageManager: "foo"})
	if err == nil {
		t.Fatal("Run() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "unsupported package manager") {
		t.Fatalf("error = %q, want unsupported package manager", err.Error())
	}
}
