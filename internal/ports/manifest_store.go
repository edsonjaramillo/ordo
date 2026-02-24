package ports

import "context"

type ManifestStore interface {
	RewriteCatalogReferences(ctx context.Context, targetDir string, catalogName string, packages []string) error
	RewriteCatalogReferencesExistingOnly(ctx context.Context, targetDir string, catalogName string, packages []string) error
}
