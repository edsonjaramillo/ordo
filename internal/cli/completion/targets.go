package completion

import (
	"context"

	"ordo/internal/app"
)

type TargetCompleter struct {
	discovery        app.DiscoveryService
	installCompleter app.InstallCompletionService
}

func NewTargetCompleter(discovery app.DiscoveryService, installCompleter app.InstallCompletionService) TargetCompleter {
	return TargetCompleter{discovery: discovery, installCompleter: installCompleter}
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

func (c TargetCompleter) WorkspaceKeys(ctx context.Context, prefix string) ([]string, error) {
	return c.installCompleter.WorkspaceKeys(ctx, prefix)
}

func (c TargetCompleter) InstallPackages(ctx context.Context, prefix string) ([]string, error) {
	return c.installCompleter.PackageSpecs(ctx, prefix)
}
