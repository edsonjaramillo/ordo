package app

import (
	"context"
	"fmt"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type GlobalUninstallRequest struct {
	Manager  domain.PackageManager
	Packages []string
}

type GlobalUninstallUseCase struct {
	runner ports.Runner
	lister ports.GlobalPackageLister
}

func NewGlobalUninstallUseCase(runner ports.Runner, lister ports.GlobalPackageLister) GlobalUninstallUseCase {
	return GlobalUninstallUseCase{runner: runner, lister: lister}
}

func (u GlobalUninstallUseCase) Run(ctx context.Context, req GlobalUninstallRequest) error {
	if u.lister == nil {
		return fmt.Errorf("global package lister is not configured")
	}

	pkgs := trimNonEmpty(req.Packages)
	installed, err := u.lister.ListInstalledGlobalPackages(ctx, req.Manager)
	if err != nil {
		return err
	}

	installedSet := map[string]struct{}{}
	for _, item := range installed {
		installedSet[item] = struct{}{}
	}

	missing := make([]string, 0, len(pkgs))
	for _, pkg := range pkgs {
		if _, ok := installedSet[pkg]; ok {
			continue
		}
		missing = append(missing, pkg)
	}
	if len(missing) > 0 {
		paths, pathErr := u.lister.ResolveGlobalStorePaths(ctx, req.Manager)
		if pathErr != nil {
			return pathErr
		}
		return GlobalPackageMissingError{
			Manager:      req.Manager,
			Missing:      missing,
			CheckedPaths: paths,
		}
	}

	argv, err := domain.BuildGlobalUninstallCommand(req.Manager, pkgs)
	if err != nil {
		return err
	}

	return u.runner.Run(ctx, ".", argv)
}
