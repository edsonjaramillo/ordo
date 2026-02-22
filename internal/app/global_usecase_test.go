package app

import (
	"context"
	"errors"
	"testing"

	"ordo/internal/domain"
)

type fakeGlobalLister struct {
	items []string
	paths []string
	err   error
}

func (f fakeGlobalLister) ListInstalledGlobalPackages(context.Context, domain.PackageManager) ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return append([]string(nil), f.items...), nil
}

func (f fakeGlobalLister) ResolveGlobalStorePaths(context.Context, domain.PackageManager) ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return append([]string(nil), f.paths...), nil
}

func TestGlobalInstallUseCase(t *testing.T) {
	runner := &fakeRunner{}
	uc := NewGlobalInstallUseCase(runner)

	err := uc.Run(context.Background(), GlobalInstallRequest{
		Manager:  domain.ManagerPNPM,
		Packages: []string{"typescript"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if runner.dir != "." {
		t.Fatalf("expected root dir '.', got %s", runner.dir)
	}
	want := []string{"pnpm", "add", "--global", "typescript"}
	if len(runner.argv) != len(want) {
		t.Fatalf("unexpected argv len: %#v", runner.argv)
	}
	for i := range want {
		if runner.argv[i] != want[i] {
			t.Fatalf("argv[%d] = %q, want %q", i, runner.argv[i], want[i])
		}
	}
}

func TestGlobalUninstallUseCase(t *testing.T) {
	runner := &fakeRunner{}
	lister := fakeGlobalLister{items: []string{"typescript"}}
	uc := NewGlobalUninstallUseCase(runner, lister)

	err := uc.Run(context.Background(), GlobalUninstallRequest{
		Manager:  domain.ManagerPNPM,
		Packages: []string{"typescript"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if runner.dir != "." {
		t.Fatalf("expected root dir '.', got %s", runner.dir)
	}
	want := []string{"pnpm", "remove", "--global", "typescript"}
	if len(runner.argv) != len(want) {
		t.Fatalf("unexpected argv len: %#v", runner.argv)
	}
	for i := range want {
		if runner.argv[i] != want[i] {
			t.Fatalf("argv[%d] = %q, want %q", i, runner.argv[i], want[i])
		}
	}
}

func TestGlobalUninstallUseCasePackageMissing(t *testing.T) {
	runner := &fakeRunner{}
	lister := fakeGlobalLister{
		items: []string{"eslint"},
		paths: []string{"/tmp/pnpm/global/5/node_modules"},
	}
	uc := NewGlobalUninstallUseCase(runner, lister)

	err := uc.Run(context.Background(), GlobalUninstallRequest{
		Manager:  domain.ManagerPNPM,
		Packages: []string{"typescript"},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var missingErr GlobalPackageMissingError
	if !errors.As(err, &missingErr) {
		t.Fatalf("expected GlobalPackageMissingError, got %T", err)
	}
	if missingErr.Manager != domain.ManagerPNPM {
		t.Fatalf("manager = %s, want %s", missingErr.Manager, domain.ManagerPNPM)
	}
	if len(missingErr.Missing) != 1 || missingErr.Missing[0] != "typescript" {
		t.Fatalf("missing = %#v", missingErr.Missing)
	}
	if len(missingErr.CheckedPaths) != 1 || missingErr.CheckedPaths[0] != "/tmp/pnpm/global/5/node_modules" {
		t.Fatalf("checked paths = %#v", missingErr.CheckedPaths)
	}
	if len(runner.argv) != 0 {
		t.Fatalf("expected uninstall command not to run, got %#v", runner.argv)
	}
}

func TestGlobalCompletionServiceInstalledGlobalPackages(t *testing.T) {
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	installCompletion := NewInstallCompletionService(discovery, nil)
	svc := NewGlobalCompletionService(installCompletion, fakeGlobalLister{items: []string{"eslint", "typescript", "eslint"}})

	items, err := svc.InstalledGlobalPackages(context.Background(), domain.ManagerPNPM, "e")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0] != "eslint" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestGlobalCompletionServiceInstalledGlobalPackagesFallback(t *testing.T) {
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	installCompletion := NewInstallCompletionService(discovery, nil)
	svc := NewGlobalCompletionService(installCompletion, fakeGlobalLister{err: errors.New("boom")})

	items, err := svc.InstalledGlobalPackages(context.Background(), domain.ManagerPNPM, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty items, got %#v", items)
	}
}

func TestGlobalCompletionServiceInstallPackageSpecs(t *testing.T) {
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	installCompletion := NewInstallCompletionService(discovery, fakeSuggestor{items: []string{"react-dom", "react"}})
	svc := NewGlobalCompletionService(installCompletion, fakeGlobalLister{})

	items, err := svc.InstallPackageSpecs(context.Background(), "react")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 || items[0] != "react" || items[1] != "react-dom" {
		t.Fatalf("unexpected items: %#v", items)
	}
}
