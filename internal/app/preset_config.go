package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"ordo/internal/config"
	"ordo/internal/domain"
	"ordo/internal/ports"
)

type ordoConfig struct {
	Presets map[string]presetConfig `json:"presets"`
}

type presetConfig struct {
	Dependencies         []string `json:"dependencies"`
	DevDependencies      []string `json:"devDependencies"`
	PeerDependencies     []string `json:"peerDependencies"`
	OptionalDependencies []string `json:"optionalDependencies"`
}

type presetConfigService struct {
	configStore ports.ConfigStore
}

func newPresetConfigService(configStore ports.ConfigStore) presetConfigService {
	return presetConfigService{configStore: configStore}
}

func (s presetConfigService) load() (ordoConfig, error) {
	path, err := config.OrdoConfigPath()
	if err != nil {
		return ordoConfig{}, err
	}

	payload, err := s.configStore.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ordoConfig{}, ErrConfigNotFound
		}
		return ordoConfig{}, err
	}

	var cfg ordoConfig
	if err := json.Unmarshal(payload, &cfg); err != nil {
		return ordoConfig{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if cfg.Presets == nil {
		cfg.Presets = map[string]presetConfig{}
	}
	return cfg, nil
}

func (s presetConfigService) presetNames(prefix string) ([]string, error) {
	cfg, err := s.load()
	if err != nil {
		return nil, err
	}

	items := make([]string, 0, len(cfg.Presets))
	for name := range cfg.Presets {
		if prefix != "" && !strings.HasPrefix(name, prefix) {
			continue
		}
		items = append(items, name)
	}
	sort.Strings(items)
	return items, nil
}

func (s presetConfigService) bucketPackages(name string, bucket domain.PresetBucket) ([]string, error) {
	cfg, err := s.load()
	if err != nil {
		return nil, err
	}

	preset, ok := cfg.Presets[strings.TrimSpace(name)]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrPresetNotFound, name)
	}

	var packages []string
	switch bucket {
	case domain.BucketDependencies:
		packages = preset.Dependencies
	case domain.BucketDevDependencies:
		packages = preset.DevDependencies
	case domain.BucketPeerDependencies:
		packages = preset.PeerDependencies
	case domain.BucketOptionalDependencies:
		packages = preset.OptionalDependencies
	default:
		packages = nil
	}

	if len(packages) == 0 {
		return nil, fmt.Errorf("%w: %s/%s", ErrPresetBucketNotFound, name, bucket)
	}
	return trimUnique(packages), nil
}

func trimUnique(items []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}
