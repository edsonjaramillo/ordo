package app

import (
	"context"
	"strings"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type InstallRequest struct {
	Packages  []string
	Workspace string
	Dev       bool
	Peer      bool
	Optional  bool
	Prod      bool
	Exact     bool
}

type InstallUseCase struct {
	discovery DiscoveryService
	runner    ports.Runner
}

func NewInstallUseCase(discovery DiscoveryService, runner ports.Runner) InstallUseCase {
	return InstallUseCase{discovery: discovery, runner: runner}
}

func (u InstallUseCase) Run(ctx context.Context, req InstallRequest) error {
	snapshot, err := u.discovery.Snapshot(ctx)
	if err != nil {
		return err
	}

	pkg, err := resolveInstallTargetPackage(snapshot, req.Workspace)
	if err != nil {
		return err
	}

	argv, err := domain.BuildInstallCommand(snapshot.Manager, trimNonEmpty(req.Packages), domain.InstallOptions{
		Dev:      req.Dev,
		Peer:     req.Peer,
		Optional: req.Optional,
		Prod:     req.Prod,
		Exact:    req.Exact,
	})
	if err != nil {
		return err
	}

	return u.runner.Run(ctx, pkg.Dir, argv)
}

func trimNonEmpty(items []string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}
