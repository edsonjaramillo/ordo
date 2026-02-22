package app

import (
	"context"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type GlobalInstallRequest struct {
	Manager  domain.PackageManager
	Packages []string
}

type GlobalInstallUseCase struct {
	runner ports.Runner
}

func NewGlobalInstallUseCase(runner ports.Runner) GlobalInstallUseCase {
	return GlobalInstallUseCase{runner: runner}
}

func (u GlobalInstallUseCase) Run(ctx context.Context, req GlobalInstallRequest) error {
	argv, err := domain.BuildGlobalInstallCommand(req.Manager, trimNonEmpty(req.Packages))
	if err != nil {
		return err
	}

	return u.runner.Run(ctx, ".", argv)
}
