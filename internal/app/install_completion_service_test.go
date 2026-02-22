package app

import (
	"context"
	"errors"
	"testing"
)

type fakeSuggestor struct {
	items []string
	err   error
}

func (f fakeSuggestor) Suggest(context.Context, string, int) ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return append([]string(nil), f.items...), nil
}

func TestInstallCompletionServiceWorkspaceKeys(t *testing.T) {
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	svc := NewInstallCompletionService(discovery, nil)

	items, err := svc.WorkspaceKeys(context.Background(), "u")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0] != "ui" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestInstallCompletionServicePackageSpecsMerged(t *testing.T) {
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	svc := NewInstallCompletionService(discovery, fakeSuggestor{items: []string{"react-dom", "react"}})

	items, err := svc.PackageSpecs(context.Background(), "react")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 || items[0] != "react" || items[1] != "react-dom" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestInstallCompletionServicePackageSpecsFallback(t *testing.T) {
	discovery := NewDiscoveryService(fakeIndexer{infos: fixtureInfos()})
	svc := NewInstallCompletionService(discovery, fakeSuggestor{err: errors.New("network down")})

	items, err := svc.PackageSpecs(context.Background(), "react")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0] != "react" {
		t.Fatalf("unexpected items: %#v", items)
	}
}
