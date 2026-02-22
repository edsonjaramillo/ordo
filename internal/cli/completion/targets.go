package completion

import (
	"context"

	"ordo/internal/app"
)

type TargetCompleter struct {
	discovery app.DiscoveryService
}

func NewTargetCompleter(discovery app.DiscoveryService) TargetCompleter {
	return TargetCompleter{discovery: discovery}
}

func (c TargetCompleter) ScriptTargets(ctx context.Context, prefix string) ([]string, error) {
	snapshot, err := c.discovery.Snapshot(ctx)
	if err != nil {
		return nil, err
	}
	return snapshot.ScriptTargets(prefix), nil
}

func (c TargetCompleter) PackageTargets(ctx context.Context, prefix string) ([]string, error) {
	snapshot, err := c.discovery.Snapshot(ctx)
	if err != nil {
		return nil, err
	}
	return snapshot.PackageTargets(prefix), nil
}
