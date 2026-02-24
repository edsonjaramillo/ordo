package app

import (
	"context"
	"errors"
	"testing"

	"ordo/internal/domain"
)

type fakeCatalogStore struct {
	manager       domain.PackageManager
	name          string
	entries       map[string]string
	removed       []string
	force         bool
	err           error
	named         []string
	catalogByName map[string][]string
}

func (f *fakeCatalogStore) UpsertCatalogEntries(_ context.Context, manager domain.PackageManager, name string, entries map[string]string, force bool) error {
	f.manager = manager
	f.name = name
	f.entries = entries
	f.force = force
	return f.err
}

func (f *fakeCatalogStore) RemoveCatalogEntries(_ context.Context, manager domain.PackageManager, name string, packages []string) error {
	f.manager = manager
	f.name = name
	f.removed = append([]string(nil), packages...)
	return f.err
}

func (f *fakeCatalogStore) NamedCatalogs(context.Context, domain.PackageManager) ([]string, error) {
	return f.named, f.err
}

func (f *fakeCatalogStore) CatalogPackageNames(_ context.Context, _ domain.PackageManager, name string) ([]string, error) {
	if f.catalogByName == nil {
		return []string{}, f.err
	}
	return f.catalogByName[name], f.err
}

type fakeManifestStore struct {
	dir      string
	name     string
	packages []string
	sync     []manifestRewriteCall
	err      error
}

type manifestRewriteCall struct {
	dir      string
	name     string
	packages []string
}

func (f *fakeManifestStore) RewriteCatalogReferences(_ context.Context, targetDir string, catalogName string, packages []string) error {
	f.dir = targetDir
	f.name = catalogName
	f.packages = append([]string(nil), packages...)
	return f.err
}

func (f *fakeManifestStore) RewriteCatalogReferencesExistingOnly(_ context.Context, targetDir string, catalogName string, packages []string) error {
	f.sync = append(f.sync, manifestRewriteCall{
		dir:      targetDir,
		name:     catalogName,
		packages: append([]string(nil), packages...),
	})
	return f.err
}

type fakeVersionResolver struct {
	versions map[string]string
	err      error
}

func (f fakeVersionResolver) LatestVersion(_ context.Context, packageName string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	if version, ok := f.versions[packageName]; ok {
		return version, nil
	}
	return "", errors.New("missing")
}

func TestCatalogUseCaseDefaultCatalogAdd(t *testing.T) {
	catalogs := &fakeCatalogStore{}
	manifests := &fakeManifestStore{}
	resolver := fakeVersionResolver{versions: map[string]string{"react": "19.1.0"}}
	uc := NewCatalogUseCase(NewDiscoveryService(fakeIndexer{infos: fixtureInfos()}), catalogs, manifests, resolver)

	err := uc.RunAdd(context.Background(), CatalogAddRequest{Packages: []string{"react"}, Force: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if catalogs.manager != domain.ManagerPNPM {
		t.Fatalf("manager = %s, want %s", catalogs.manager, domain.ManagerPNPM)
	}
	if catalogs.name != "" {
		t.Fatalf("name = %q, want default", catalogs.name)
	}
	if catalogs.entries["react"] != "19.1.0" {
		t.Fatalf("entries = %#v", catalogs.entries)
	}
	if !catalogs.force {
		t.Fatal("force flag not propagated")
	}
	if manifests.dir != "." || manifests.name != "" {
		t.Fatalf("manifest target = %q name=%q", manifests.dir, manifests.name)
	}
	if len(manifests.packages) != 1 || manifests.packages[0] != "react" {
		t.Fatalf("manifest packages = %#v", manifests.packages)
	}
}

func TestCatalogUseCaseNamedCatalogAddWorkspace(t *testing.T) {
	catalogs := &fakeCatalogStore{}
	manifests := &fakeManifestStore{}
	resolver := fakeVersionResolver{}
	uc := NewCatalogUseCase(NewDiscoveryService(fakeIndexer{infos: fixtureInfos()}), catalogs, manifests, resolver)

	err := uc.RunNamedAdd(context.Background(), CatalogsAddRequest{Name: "react19", Workspace: "ui", Packages: []string{"react@19.1.0"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if manifests.dir != "packages/ui" {
		t.Fatalf("manifest dir = %q, want packages/ui", manifests.dir)
	}
	if manifests.name != "react19" {
		t.Fatalf("manifest name = %q", manifests.name)
	}
	if catalogs.entries["react"] != "19.1.0" {
		t.Fatalf("entries = %#v", catalogs.entries)
	}
}

func TestCatalogUseCaseUnsupportedManager(t *testing.T) {
	catalogs := &fakeCatalogStore{}
	manifests := &fakeManifestStore{}
	resolver := fakeVersionResolver{}
	discovery := NewDiscoveryService(fakeIndexer{infos: []domain.PackageInfo{{Dir: ".", Lockfiles: map[string]bool{"package-lock.json": true}}}})
	uc := NewCatalogUseCase(discovery, catalogs, manifests, resolver)

	err := uc.RunAdd(context.Background(), CatalogAddRequest{Packages: []string{"react@1.0.0"}})
	if !errors.Is(err, ErrCatalogUnsupported) {
		t.Fatalf("expected ErrCatalogUnsupported, got %v", err)
	}
}

func TestCatalogUseCaseConflict(t *testing.T) {
	catalogs := &fakeCatalogStore{err: errors.New("catalog conflict")}
	manifests := &fakeManifestStore{}
	resolver := fakeVersionResolver{}
	uc := NewCatalogUseCase(NewDiscoveryService(fakeIndexer{infos: fixtureInfos()}), catalogs, manifests, resolver)

	err := uc.RunAdd(context.Background(), CatalogAddRequest{Packages: []string{"react@1.0.0"}})
	if !errors.Is(err, ErrCatalogConflict) {
		t.Fatalf("expected ErrCatalogConflict, got %v", err)
	}
}

func TestCatalogUseCaseInvalidName(t *testing.T) {
	uc := NewCatalogUseCase(NewDiscoveryService(fakeIndexer{infos: fixtureInfos()}), &fakeCatalogStore{}, &fakeManifestStore{}, fakeVersionResolver{})

	err := uc.RunNamedAdd(context.Background(), CatalogsAddRequest{Name: "React19", Packages: []string{"react@19.1.0"}})
	if !errors.Is(err, ErrInvalidCatalogName) {
		t.Fatalf("expected ErrInvalidCatalogName, got %v", err)
	}
}

func TestCatalogUseCaseRemove(t *testing.T) {
	catalogs := &fakeCatalogStore{}
	manifests := &fakeManifestStore{}
	uc := NewCatalogUseCase(NewDiscoveryService(fakeIndexer{infos: fixtureInfos()}), catalogs, manifests, fakeVersionResolver{})

	err := uc.RunRemove(context.Background(), CatalogRemoveRequest{Packages: []string{"react", "react"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(catalogs.removed) != 1 || catalogs.removed[0] != "react" {
		t.Fatalf("removed = %#v", catalogs.removed)
	}
	if manifests.dir != "" {
		t.Fatalf("manifest should not be touched, got dir=%q", manifests.dir)
	}
}

func TestCatalogUseCaseRemoveRejectsVersion(t *testing.T) {
	uc := NewCatalogUseCase(NewDiscoveryService(fakeIndexer{infos: fixtureInfos()}), &fakeCatalogStore{}, &fakeManifestStore{}, fakeVersionResolver{})

	err := uc.RunRemove(context.Background(), CatalogRemoveRequest{Packages: []string{"react@19"}})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCatalogUseCaseSync(t *testing.T) {
	catalogs := &fakeCatalogStore{
		catalogByName: map[string][]string{
			"": []string{"react", "zod"},
		},
	}
	manifests := &fakeManifestStore{}
	uc := NewCatalogUseCase(NewDiscoveryService(fakeIndexer{infos: fixtureInfos()}), catalogs, manifests, fakeVersionResolver{})

	err := uc.RunSync(context.Background(), CatalogSyncRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(manifests.sync) != 1 {
		t.Fatalf("sync calls = %#v", manifests.sync)
	}
	if manifests.sync[0].dir != "packages/ui" {
		t.Fatalf("sync dir = %q, want packages/ui", manifests.sync[0].dir)
	}
	if manifests.sync[0].name != "" {
		t.Fatalf("sync name = %q, want default", manifests.sync[0].name)
	}
	if len(manifests.sync[0].packages) != 2 || manifests.sync[0].packages[0] != "react" || manifests.sync[0].packages[1] != "zod" {
		t.Fatalf("sync packages = %#v", manifests.sync[0].packages)
	}
}

func TestCatalogUseCaseSyncNoCatalogEntries(t *testing.T) {
	catalogs := &fakeCatalogStore{catalogByName: map[string][]string{"": []string{}}}
	manifests := &fakeManifestStore{}
	uc := NewCatalogUseCase(NewDiscoveryService(fakeIndexer{infos: fixtureInfos()}), catalogs, manifests, fakeVersionResolver{})

	err := uc.RunSync(context.Background(), CatalogSyncRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(manifests.sync) != 0 {
		t.Fatalf("sync calls = %#v, want none", manifests.sync)
	}
}

func TestCatalogUseCaseSyncNoWorkspaces(t *testing.T) {
	catalogs := &fakeCatalogStore{catalogByName: map[string][]string{"": {"react"}}}
	manifests := &fakeManifestStore{}
	infos := []domain.PackageInfo{
		{Dir: ".", Lockfiles: map[string]bool{"pnpm-lock.yaml": true}},
	}
	uc := NewCatalogUseCase(NewDiscoveryService(fakeIndexer{infos: infos}), catalogs, manifests, fakeVersionResolver{})

	err := uc.RunSync(context.Background(), CatalogSyncRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(manifests.sync) != 0 {
		t.Fatalf("sync calls = %#v, want none", manifests.sync)
	}
}

func TestCatalogUseCaseSyncUnsupportedManager(t *testing.T) {
	catalogs := &fakeCatalogStore{catalogByName: map[string][]string{"": {"react"}}}
	manifests := &fakeManifestStore{}
	infos := []domain.PackageInfo{
		{Dir: ".", Lockfiles: map[string]bool{"package-lock.json": true}},
		{Dir: "packages/ui"},
	}
	uc := NewCatalogUseCase(NewDiscoveryService(fakeIndexer{infos: infos}), catalogs, manifests, fakeVersionResolver{})

	err := uc.RunSync(context.Background(), CatalogSyncRequest{})
	if !errors.Is(err, ErrCatalogUnsupported) {
		t.Fatalf("expected ErrCatalogUnsupported, got %v", err)
	}
}
