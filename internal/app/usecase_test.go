package app

import (
	"context"
	"errors"
	"testing"

	"ordo/internal/domain"
)

type fakeIndexer struct {
	infos []domain.PackageInfo
	err   error
}

func (f fakeIndexer) Discover(context.Context) ([]domain.PackageInfo, error) {
	return f.infos, f.err
}

type fakeRunner struct {
	dir  string
	argv []string
	err  error
}

func (f *fakeRunner) Run(_ context.Context, dir string, argv []string) error {
	f.dir = dir
	f.argv = append([]string(nil), argv...)
	return f.err
}

func TestRunUseCaseWorkspace(t *testing.T) {
	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewRunUseCase(discovery, runner)

	err := uc.Run(context.Background(), RunRequest{Target: "ui/build", ExtraArgs: []string{"--watch"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if runner.dir != "packages/ui" {
		t.Fatalf("expected runner dir packages/ui, got %s", runner.dir)
	}
	if len(runner.argv) != 4 || runner.argv[0] != "pnpm" || runner.argv[1] != "run" || runner.argv[2] != "build" || runner.argv[3] != "--watch" {
		t.Fatalf("unexpected argv: %#v", runner.argv)
	}
}

func TestRunUseCaseScriptMissing(t *testing.T) {
	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewRunUseCase(discovery, runner)

	err := uc.Run(context.Background(), RunRequest{Target: "ui/does-not-exist"})
	if !errors.Is(err, ErrScriptNotFound) {
		t.Fatalf("expected ErrScriptNotFound, got %v", err)
	}
}

func TestUninstallUseCaseRoot(t *testing.T) {
	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewUninstallUseCase(discovery, runner)

	err := uc.Run(context.Background(), UninstallRequest{Target: "react"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if runner.dir != "." {
		t.Fatalf("expected root dir '.', got %s", runner.dir)
	}
	if len(runner.argv) != 3 || runner.argv[0] != "pnpm" || runner.argv[1] != "remove" || runner.argv[2] != "react" {
		t.Fatalf("unexpected argv: %#v", runner.argv)
	}
}

func TestUninstallUseCasePackageMissing(t *testing.T) {
	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewUninstallUseCase(discovery, runner)

	err := uc.Run(context.Background(), UninstallRequest{Target: "ui/left-pad"})
	if !errors.Is(err, ErrPackageNotFound) {
		t.Fatalf("expected ErrPackageNotFound, got %v", err)
	}
}

func TestUpdateUseCaseRoot(t *testing.T) {
	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewUpdateUseCase(discovery, runner)

	err := uc.Run(context.Background(), UpdateRequest{Target: "react"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if runner.dir != "." {
		t.Fatalf("expected root dir '.', got %s", runner.dir)
	}
	if len(runner.argv) != 3 || runner.argv[0] != "pnpm" || runner.argv[1] != "update" || runner.argv[2] != "react" {
		t.Fatalf("unexpected argv: %#v", runner.argv)
	}
}

func TestUpdateUseCasePackageMissing(t *testing.T) {
	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewUpdateUseCase(discovery, runner)

	err := uc.Run(context.Background(), UpdateRequest{Target: "ui/left-pad"})
	if !errors.Is(err, ErrPackageNotFound) {
		t.Fatalf("expected ErrPackageNotFound, got %v", err)
	}
}

func TestSnapshotWorkspaceCollision(t *testing.T) {
	discovery := NewDiscoveryService(fakeIndexer{infos: []domain.PackageInfo{
		{Dir: "."},
		{Dir: "apps/web", Scripts: map[string]string{"dev": "vite"}},
		{Dir: "packages/web", Scripts: map[string]string{"build": "tsup"}},
	}})

	snapshot, err := discovery.Snapshot(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	targets := snapshot.ScriptTargets("")
	foundAppsWeb := false
	foundPackagesWeb := false
	for _, target := range targets {
		if target == "apps/web/dev" {
			foundAppsWeb = true
		}
		if target == "packages/web/build" {
			foundPackagesWeb = true
		}
	}
	if !foundAppsWeb || !foundPackagesWeb {
		t.Fatalf("expected collision fallback targets, got %#v", targets)
	}
}

func fixtureInfos() []domain.PackageInfo {
	return []domain.PackageInfo{
		{
			Dir:          ".",
			Scripts:      map[string]string{"build": "turbo run build"},
			Dependencies: map[string]struct{}{"react": {}},
			Lockfiles:    map[string]bool{"pnpm-lock.yaml": true},
		},
		{
			Dir:          "packages/ui",
			Scripts:      map[string]string{"build": "tsup"},
			Dependencies: map[string]struct{}{"clsx": {}},
		},
	}
}
