package ports

import (
	"context"

	"ordo/internal/domain"
)

type WorkspaceIndexer interface {
	Discover(ctx context.Context) ([]domain.PackageInfo, error)
}
