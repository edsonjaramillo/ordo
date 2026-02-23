package app

import (
	"context"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type GlobalUpdateRequest struct {
	Manager  domain.PackageManager
	Packages []string
}

type GlobalUpdateUseCase struct {
	runner ports.Runner
}

func NewGlobalUpdateUseCase(runner ports.Runner) GlobalUpdateUseCase {
	return GlobalUpdateUseCase{runner: runner}
}

func (u GlobalUpdateUseCase) Run(ctx context.Context, req GlobalUpdateRequest) error {
	argv, err := domain.BuildGlobalUpdateCommand(req.Manager, trimNonEmpty(req.Packages))
	if err != nil {
		return err
	}

	return u.runner.Run(ctx, ".", argv)
}
