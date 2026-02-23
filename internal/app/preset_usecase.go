package app

import (
	"context"
	"fmt"
	"strings"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type PresetRequest struct {
	Preset    string
	Bucket    string
	Packages  []string
	Workspace string
}

type PresetUseCase struct {
	discovery DiscoveryService
	runner    ports.Runner
	config    presetConfigService
}

func NewPresetUseCase(
	discovery DiscoveryService,
	runner ports.Runner,
	configStore ports.ConfigStore,
) PresetUseCase {
	return PresetUseCase{
		discovery: discovery,
		runner:    runner,
		config:    newPresetConfigService(configStore),
	}
}

func (u PresetUseCase) Run(ctx context.Context, req PresetRequest) error {
	snapshot, err := u.discovery.Snapshot(ctx)
	if err != nil {
		return err
	}

	target, err := resolveInstallTargetPackage(snapshot, req.Workspace)
	if err != nil {
		return err
	}

	bucket, err := domain.ParsePresetBucket(req.Bucket)
	if err != nil {
		return err
	}

	packages, err := u.config.bucketPackages(req.Preset, bucket)
	if err != nil {
		return err
	}

	selected, err := filterPresetPackages(packages, req.Packages)
	if err != nil {
		return err
	}

	argv, err := domain.BuildInstallCommand(snapshot.Manager, selected, domain.BucketInstallOptions(bucket))
	if err != nil {
		return err
	}

	return u.runner.Run(ctx, target.Dir, argv)
}

func filterPresetPackages(available []string, requested []string) ([]string, error) {
	availableSet := map[string]struct{}{}
	for _, item := range available {
		availableSet[item] = struct{}{}
	}

	if len(requested) == 0 {
		return available, nil
	}

	selected := make([]string, 0, len(requested))
	seen := map[string]struct{}{}
	for _, raw := range requested {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		if _, ok := availableSet[name]; !ok {
			return nil, fmt.Errorf("%w: %s", ErrPresetPackageNotFound, name)
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		selected = append(selected, name)
	}

	if len(selected) == 0 {
		return nil, fmt.Errorf("%w: no package filters provided", ErrPresetPackageNotFound)
	}
	return selected, nil
}
