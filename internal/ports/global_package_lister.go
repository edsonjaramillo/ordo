package ports

import (
	"context"

	"ordo/internal/domain"
)

type GlobalPackageLister interface {
	ListInstalledGlobalPackages(ctx context.Context, manager domain.PackageManager) ([]string, error)
	ResolveGlobalStorePaths(ctx context.Context, manager domain.PackageManager) ([]string, error)
}
