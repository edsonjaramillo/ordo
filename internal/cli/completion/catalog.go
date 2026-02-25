package completion

import (
	"context"

	"ordo/internal/app"
)

type CatalogCompleter struct {
	catalogs app.CatalogCompletionService
}

func NewCatalogCompleter(catalogs app.CatalogCompletionService) CatalogCompleter {
	return CatalogCompleter{catalogs: catalogs}
}

func (c CatalogCompleter) PackageSpecs(ctx context.Context, prefix string) ([]string, error) {
	return c.catalogs.PackageSpecs(ctx, prefix)
}

func (c CatalogCompleter) WorkspaceKeys(ctx context.Context, prefix string) ([]string, error) {
	return c.catalogs.WorkspaceKeys(ctx, prefix)
}

func (c CatalogCompleter) NamedCatalogs(ctx context.Context, prefix string) ([]string, error) {
	return c.catalogs.NamedCatalogs(ctx, prefix)
}

func (c CatalogCompleter) CatalogPackageNames(ctx context.Context, name string, prefix string) ([]string, error) {
	return c.catalogs.CatalogPackageNames(ctx, name, prefix)
}

func (c CatalogCompleter) WorkspaceDependencyNames(ctx context.Context, workspace string, prefix string) ([]string, error) {
	return c.catalogs.WorkspaceDependencyNames(ctx, workspace, prefix)
}
