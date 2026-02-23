package ports

import "context"

type PackageManagerAvailability interface {
	AvailablePackageManagers(ctx context.Context) ([]string, error)
}
