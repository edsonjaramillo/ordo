package ports

import (
	"context"

	"ordo/internal/domain"
)

type CatalogStore interface {
	UpsertCatalogEntries(ctx context.Context, manager domain.PackageManager, name string, entries map[string]string, force bool) error
	RemoveCatalogEntries(ctx context.Context, manager domain.PackageManager, name string, packages []string) error
	NamedCatalogs(ctx context.Context, manager domain.PackageManager) ([]string, error)
	CatalogPackageNames(ctx context.Context, manager domain.PackageManager, name string) ([]string, error)
}
