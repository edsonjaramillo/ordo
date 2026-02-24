package app

import (
	"context"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type CatalogCompletionService struct {
	discovery DiscoveryService
	install   InstallCompletionService
	catalogs  ports.CatalogStore
}

func NewCatalogCompletionService(
	discovery DiscoveryService,
	install InstallCompletionService,
	catalogs ports.CatalogStore,
) CatalogCompletionService {
	return CatalogCompletionService{
		discovery: discovery,
		install:   install,
		catalogs:  catalogs,
	}
}

func (s CatalogCompletionService) PackageSpecs(ctx context.Context, prefix string) ([]string, error) {
	return s.install.PackageSpecs(ctx, prefix)
}

func (s CatalogCompletionService) WorkspaceKeys(ctx context.Context, prefix string) ([]string, error) {
	return s.install.WorkspaceKeys(ctx, prefix)
}

func (s CatalogCompletionService) NamedCatalogs(ctx context.Context, prefix string) ([]string, error) {
	snapshot, err := s.discovery.Snapshot(ctx)
	if err != nil {
		return nil, err
	}
	if !domain.SupportsCatalogs(snapshot.Manager) {
		return []string{}, nil
	}
	items, err := s.catalogs.NamedCatalogs(ctx, snapshot.Manager)
	if err != nil {
		return nil, err
	}
	return filterAndSort(items, prefix), nil
}

func (s CatalogCompletionService) CatalogPackageNames(ctx context.Context, name string, prefix string) ([]string, error) {
	snapshot, err := s.discovery.Snapshot(ctx)
	if err != nil {
		return nil, err
	}
	if !domain.SupportsCatalogs(snapshot.Manager) {
		return []string{}, nil
	}
	items, err := s.catalogs.CatalogPackageNames(ctx, snapshot.Manager, name)
	if err != nil {
		return nil, err
	}
	return filterAndSort(items, prefix), nil
}
