package ports

import "context"

type PackageSuggestor interface {
	Suggest(ctx context.Context, prefix string, limit int) ([]string, error)
}
