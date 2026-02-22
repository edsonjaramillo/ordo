package app

import (
	"context"
	"fmt"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type RunRequest struct {
	Target    string
	ExtraArgs []string
}

type RunUseCase struct {
	discovery DiscoveryService
	runner    ports.Runner
}

func NewRunUseCase(discovery DiscoveryService, runner ports.Runner) RunUseCase {
	return RunUseCase{discovery: discovery, runner: runner}
}

func (u RunUseCase) Run(ctx context.Context, req RunRequest) error {
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

	if _, ok := pkg.Scripts[target.Name]; !ok {
		return fmt.Errorf("%w: %s", ErrScriptNotFound, target.Name)
	}

	argv, err := domain.BuildRunCommand(snapshot.Manager, target.Name, req.ExtraArgs)
	if err != nil {
		return err
	}

	return u.runner.Run(ctx, pkg.Dir, argv)
}
