package ports

import "context"

type PackageVersionResolver interface {
	LatestVersion(ctx context.Context, packageName string) (string, error)
}
