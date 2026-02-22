package app

import (
	"context"
	"fmt"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type UninstallRequest struct {
	Target string
}

type UninstallUseCase struct {
	discovery DiscoveryService
	runner    ports.Runner
}

func NewUninstallUseCase(discovery DiscoveryService, runner ports.Runner) UninstallUseCase {
	return UninstallUseCase{discovery: discovery, runner: runner}
}

func (u UninstallUseCase) Run(ctx context.Context, req UninstallRequest) error {
	target, err := domain.ParseTarget(req.Target)
	if err != nil {
		return err
	}

	snapshot, err := u.discovery.Snapshot(ctx)
	if err != nil {
		return err
	}

	pkg, err := resolveTargetPackage(snapshot, target)
	if err != nil {
		return err
	}

	if _, ok := pkg.Dependencies[target.Name]; !ok {
		return fmt.Errorf("%w: %s", ErrPackageNotFound, target.Name)
	}

	argv, err := domain.BuildUninstallCommand(snapshot.Manager, target.Name)
	if err != nil {
		return err
	}

	return u.runner.Run(ctx, pkg.Dir, argv)
}
