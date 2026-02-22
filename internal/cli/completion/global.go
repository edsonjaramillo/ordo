package completion

import (
	"context"

	"ordo/internal/app"
	"ordo/internal/domain"
)

type GlobalCompleter struct {
	global app.GlobalCompletionService
}

func NewGlobalCompleter(global app.GlobalCompletionService) GlobalCompleter {
	return GlobalCompleter{global: global}
}

func (c GlobalCompleter) InstallPackages(ctx context.Context, prefix string) ([]string, error) {
	return c.global.InstallPackageSpecs(ctx, prefix)
}

func (c GlobalCompleter) InstalledGlobalPackages(ctx context.Context, manager domain.PackageManager, prefix string) ([]string, error) {
	return c.global.InstalledGlobalPackages(ctx, manager, prefix)
}
