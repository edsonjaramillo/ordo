package app

import (
	"context"
	"strings"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type PresetCompletionService struct {
	config presetConfigService
}

func NewPresetCompletionService(configStore ports.ConfigStore) PresetCompletionService {
	return PresetCompletionService{config: newPresetConfigService(configStore)}
}

func (s PresetCompletionService) PresetNames(_ context.Context, prefix string) ([]string, error) {
	return s.config.presetNames(prefix)
}

func (s PresetCompletionService) Buckets(_ context.Context, preset string, prefix string) ([]string, error) {
	cfg, err := s.config.load()
	if err != nil {
		return nil, err
	}
	selected, ok := cfg.Presets[strings.TrimSpace(preset)]
	if !ok {
		return []string{}, nil
	}

	return filterAndSort(nonEmptyPresetBuckets(selected), prefix), nil
}

func (s PresetCompletionService) BucketPackages(_ context.Context, preset string, bucket string, prefix string) ([]string, error) {
	parsedBucket, err := domain.ParsePresetBucket(bucket)
	if err != nil {
		return []string{}, nil
	}

	items, err := s.config.bucketPackages(preset, parsedBucket)
	if err != nil {
		return nil, err
	}
	return filterAndSort(items, prefix), nil
}

func nonEmptyPresetBuckets(preset presetConfig) []string {
	items := make([]string, 0, 4)
	if len(trimUnique(preset.Dependencies)) > 0 {
		items = append(items, string(domain.BucketDependencies))
	}
	if len(trimUnique(preset.DevDependencies)) > 0 {
		items = append(items, string(domain.BucketDevDependencies))
	}
	if len(trimUnique(preset.PeerDependencies)) > 0 {
		items = append(items, string(domain.BucketPeerDependencies))
	}
	if len(trimUnique(preset.OptionalDependencies)) > 0 {
		items = append(items, string(domain.BucketOptionalDependencies))
	}
	return items
}
