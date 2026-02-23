package app

import (
	"context"
	"fmt"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type UpdateRequest struct {
	Target string
}

type UpdateUseCase struct {
	discovery DiscoveryService
	runner    ports.Runner
}

func NewUpdateUseCase(discovery DiscoveryService, runner ports.Runner) UpdateUseCase {
	return UpdateUseCase{discovery: discovery, runner: runner}
}

func (u UpdateUseCase) Run(ctx context.Context, req UpdateRequest) error {
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

	argv, err := domain.BuildUpdateCommand(snapshot.Manager, target.Name)
	if err != nil {
		return err
	}

	return u.runner.Run(ctx, pkg.Dir, argv)
}
