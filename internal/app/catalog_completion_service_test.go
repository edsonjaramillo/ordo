package app

import (
	"context"
	"testing"

	"ordo/internal/domain"
)

func TestCatalogCompletionServiceNamedCatalogs(t *testing.T) {
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	install := NewInstallCompletionService(discovery, nil)
	store := &fakeCatalogStore{named: []string{"react18", "react19"}}
	svc := NewCatalogCompletionService(discovery, install, store)

	items, err := svc.NamedCatalogs(context.Background(), "react1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 || items[0] != "react18" || items[1] != "react19" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestCatalogCompletionServiceNamedCatalogsUnsupportedManager(t *testing.T) {
	discovery := NewDiscoveryService(fakeIndexer{infos: []domain.PackageInfo{{Dir: ".", Lockfiles: map[string]bool{"package-lock.json": true}}}})
	install := NewInstallCompletionService(discovery, nil)
	store := &fakeCatalogStore{named: []string{"react19"}}
	svc := NewCatalogCompletionService(discovery, install, store)

	items, err := svc.NamedCatalogs(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected no items, got %#v", items)
	}
}

func TestCatalogCompletionServiceCatalogPackageNames(t *testing.T) {
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	install := NewInstallCompletionService(discovery, nil)
	store := &fakeCatalogStore{catalogByName: map[string][]string{"": {"react", "react-dom"}}}
	svc := NewCatalogCompletionService(discovery, install, store)

	items, err := svc.CatalogPackageNames(context.Background(), "", "react")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 || items[0] != "react" || items[1] != "react-dom" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestCatalogCompletionServiceWorkspaceDependencyNames(t *testing.T) {
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	install := NewInstallCompletionService(discovery, nil)
	svc := NewCatalogCompletionService(discovery, install, &fakeCatalogStore{})

	items, err := svc.WorkspaceDependencyNames(context.Background(), "ui", "cl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0] != "clsx" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestCatalogCompletionServiceWorkspaceDependencyNamesFiltersCatalogRefs(t *testing.T) {
	infos := fixtureInfos()
	infos[1].Dependencies["zod"] = struct{}{}
	infos[1].DependencyVersions["zod"] = "catalog:"
	discovery := NewDiscoveryService(fakeIndexer{infos: infos})
	install := NewInstallCompletionService(discovery, nil)
	svc := NewCatalogCompletionService(discovery, install, &fakeCatalogStore{})

	items, err := svc.WorkspaceDependencyNames(context.Background(), "ui", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0] != "clsx" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestCatalogCompletionServiceWorkspaceDependencyNamesUnknownWorkspace(t *testing.T) {
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	install := NewInstallCompletionService(discovery, nil)
	svc := NewCatalogCompletionService(discovery, install, &fakeCatalogStore{})

	items, err := svc.WorkspaceDependencyNames(context.Background(), "missing", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected no items, got %#v", items)
	}
}
