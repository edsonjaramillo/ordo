package app

import (
	"context"
	"errors"
	"testing"
)

func TestInstallUseCaseRoot(t *testing.T) {
	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewInstallUseCase(discovery, runner)

	err := uc.Run(context.Background(), InstallRequest{
		Packages: []string{"typescript"},
		Dev:      true,
		Exact:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if runner.dir != "." {
		t.Fatalf("expected root dir '.', got %s", runner.dir)
	}
	want := []string{"pnpm", "add", "--save-dev", "--save-exact", "typescript"}
	if len(runner.argv) != len(want) {
		t.Fatalf("unexpected argv len: %#v", runner.argv)
	}
	for i := range want {
		if runner.argv[i] != want[i] {
			t.Fatalf("argv[%d] = %q, want %q", i, runner.argv[i], want[i])
		}
	}
}

func TestInstallUseCaseWorkspace(t *testing.T) {
	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewInstallUseCase(discovery, runner)

	err := uc.Run(context.Background(), InstallRequest{
		Packages:  []string{"zod"},
		Workspace: "ui",
		Peer:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if runner.dir != "packages/ui" {
		t.Fatalf("expected runner dir packages/ui, got %s", runner.dir)
	}
	want := []string{"pnpm", "add", "--save-peer", "zod"}
	if len(runner.argv) != len(want) {
		t.Fatalf("unexpected argv len: %#v", runner.argv)
	}
	for i := range want {
		if runner.argv[i] != want[i] {
			t.Fatalf("argv[%d] = %q, want %q", i, runner.argv[i], want[i])
		}
	}
}

func TestInstallUseCaseWorkspaceMissing(t *testing.T) {
	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewInstallUseCase(discovery, runner)

	err := uc.Run(context.Background(), InstallRequest{
		Packages:  []string{"zod"},
		Workspace: "does-not-exist",
	})
	if !errors.Is(err, ErrWorkspaceNotFound) {
		t.Fatalf("expected ErrWorkspaceNotFound, got %v", err)
	}
}

func TestInstallUseCaseConflictingFlags(t *testing.T) {
	runner := &fakeRunner{}
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	uc := NewInstallUseCase(discovery, runner)

	err := uc.Run(context.Background(), InstallRequest{
		Packages: []string{"zod"},
		Dev:      true,
		Optional: true,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
