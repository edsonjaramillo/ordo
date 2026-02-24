package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"ordo/internal/domain"
	"ordo/internal/ports"

	"gopkg.in/yaml.v3"
)

type Store struct {
	root string
	fs   ports.ConfigStore
}

func NewStore(root string, fs ports.ConfigStore) Store {
	return Store{root: root, fs: fs}
}

func (s Store) UpsertCatalogEntries(_ context.Context, manager domain.PackageManager, name string, entries map[string]string, force bool) error {
	switch manager {
	case domain.ManagerBun:
		return s.upsertBun(name, entries, force)
	case domain.ManagerPNPM:
		return s.upsertPNPM(name, entries, force)
	case domain.ManagerYarn:
		return s.upsertYarn(name, entries, force)
	default:
		return fmt.Errorf("catalogs unsupported for package manager: %s", manager)
	}
}

func (s Store) RemoveCatalogEntries(_ context.Context, manager domain.PackageManager, name string, packages []string) error {
	switch manager {
	case domain.ManagerBun:
		return s.removeBun(name, packages)
	case domain.ManagerPNPM:
		return s.removePNPM(name, packages)
	case domain.ManagerYarn:
		return s.removeYarn(name, packages)
	default:
		return fmt.Errorf("catalogs unsupported for package manager: %s", manager)
	}
}

func (s Store) NamedCatalogs(_ context.Context, manager domain.PackageManager) ([]string, error) {
	switch manager {
	case domain.ManagerBun:
		return s.namedBunCatalogs()
	case domain.ManagerPNPM:
		return s.namedPNPMCatalogs()
	case domain.ManagerYarn:
		return s.namedYarnCatalogs()
	default:
		return []string{}, nil
	}
}

func (s Store) CatalogPackageNames(_ context.Context, manager domain.PackageManager, name string) ([]string, error) {
	switch manager {
	case domain.ManagerBun:
		return s.catalogPackageNamesBun(name)
	case domain.ManagerPNPM:
		return s.catalogPackageNamesPNPM(name)
	case domain.ManagerYarn:
		return s.catalogPackageNamesYarn(name)
	default:
		return []string{}, nil
	}
}

func (s Store) upsertBun(name string, entries map[string]string, force bool) error {
	path := filepath.Join(s.root, "package.json")
	payload := map[string]any{}

	content, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			payload = map[string]any{}
		} else {
			return err
		}
	} else if len(content) > 0 {
		if err := json.Unmarshal(content, &payload); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
	}

	if strings.TrimSpace(name) == "" {
		next, err := upsertStringMap(anyToStringMap(payload["catalog"]), entries, force)
		if err != nil {
			return err
		}
		payload["catalog"] = next
	} else {
		catalogs := anyToStringMapMap(payload["catalogs"])
		selected := catalogs[name]
		next, err := upsertStringMap(selected, entries, force)
		if err != nil {
			return err
		}
		catalogs[name] = next
		payload["catalogs"] = catalogs
	}

	formatted, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	formatted = append(formatted, '\n')
	return s.fs.WriteFile(path, formatted, 0o644)
}

func (s Store) upsertPNPM(name string, entries map[string]string, force bool) error {
	path := filepath.Join(s.root, "pnpm-workspace.yaml")
	payload, err := s.loadYAML(path)
	if err != nil {
		return err
	}

	if strings.TrimSpace(name) == "" {
		next, err := upsertStringMap(anyToStringMap(payload["catalog"]), entries, force)
		if err != nil {
			return err
		}
		payload["catalog"] = next
	} else {
		catalogs := anyToStringMapMap(payload["catalogs"])
		selected := catalogs[name]
		next, err := upsertStringMap(selected, entries, force)
		if err != nil {
			return err
		}
		catalogs[name] = next
		payload["catalogs"] = catalogs
	}

	return s.writeYAML(path, payload)
}

func (s Store) upsertYarn(name string, entries map[string]string, force bool) error {
	path := filepath.Join(s.root, ".yarnrc.yml")
	payload, err := s.loadYAML(path)
	if err != nil {
		return err
	}

	if strings.TrimSpace(name) == "" {
		next, err := upsertStringMap(anyToStringMap(payload["npmCatalog"]), entries, force)
		if err != nil {
			return err
		}
		payload["npmCatalog"] = next
	} else {
		catalogs := anyToStringMapMap(payload["npmCatalogs"])
		selected := catalogs[name]
		next, err := upsertStringMap(selected, entries, force)
		if err != nil {
			return err
		}
		catalogs[name] = next
		payload["npmCatalogs"] = catalogs
	}

	return s.writeYAML(path, payload)
}

func (s Store) removeBun(name string, packages []string) error {
	path := filepath.Join(s.root, "package.json")
	payload := map[string]any{}

	content, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(content) == 0 {
		return nil
	}
	if err := json.Unmarshal(content, &payload); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	if strings.TrimSpace(name) == "" {
		next := removeFromStringMap(anyToStringMap(payload["catalog"]), packages)
		if len(next) == 0 {
			delete(payload, "catalog")
		} else {
			payload["catalog"] = next
		}
	} else {
		catalogs := anyToStringMapMap(payload["catalogs"])
		next := removeFromStringMap(catalogs[name], packages)
		if len(next) == 0 {
			delete(catalogs, name)
		} else {
			catalogs[name] = next
		}
		if len(catalogs) == 0 {
			delete(payload, "catalogs")
		} else {
			payload["catalogs"] = catalogs
		}
	}

	formatted, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	formatted = append(formatted, '\n')
	return s.fs.WriteFile(path, formatted, 0o644)
}

func (s Store) removePNPM(name string, packages []string) error {
	path := filepath.Join(s.root, "pnpm-workspace.yaml")
	payload, err := s.loadYAML(path)
	if err != nil {
		return err
	}

	if strings.TrimSpace(name) == "" {
		next := removeFromStringMap(anyToStringMap(payload["catalog"]), packages)
		if len(next) == 0 {
			delete(payload, "catalog")
		} else {
			payload["catalog"] = next
		}
	} else {
		catalogs := anyToStringMapMap(payload["catalogs"])
		next := removeFromStringMap(catalogs[name], packages)
		if len(next) == 0 {
			delete(catalogs, name)
		} else {
			catalogs[name] = next
		}
		if len(catalogs) == 0 {
			delete(payload, "catalogs")
		} else {
			payload["catalogs"] = catalogs
		}
	}

	return s.writeYAML(path, payload)
}

func (s Store) removeYarn(name string, packages []string) error {
	path := filepath.Join(s.root, ".yarnrc.yml")
	payload, err := s.loadYAML(path)
	if err != nil {
		return err
	}

	if strings.TrimSpace(name) == "" {
		next := removeFromStringMap(anyToStringMap(payload["npmCatalog"]), packages)
		if len(next) == 0 {
			delete(payload, "npmCatalog")
		} else {
			payload["npmCatalog"] = next
		}
	} else {
		catalogs := anyToStringMapMap(payload["npmCatalogs"])
		next := removeFromStringMap(catalogs[name], packages)
		if len(next) == 0 {
			delete(catalogs, name)
		} else {
			catalogs[name] = next
		}
		if len(catalogs) == 0 {
			delete(payload, "npmCatalogs")
		} else {
			payload["npmCatalogs"] = catalogs
		}
	}

	return s.writeYAML(path, payload)
}

func (s Store) namedBunCatalogs() ([]string, error) {
	path := filepath.Join(s.root, "package.json")
	content, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}

	payload := map[string]any{}
	if err := json.Unmarshal(content, &payload); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return sortedMapKeys(anyToStringMapMap(payload["catalogs"])), nil
}

func (s Store) namedPNPMCatalogs() ([]string, error) {
	path := filepath.Join(s.root, "pnpm-workspace.yaml")
	payload, err := s.loadYAML(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	return sortedMapKeys(anyToStringMapMap(payload["catalogs"])), nil
}

func (s Store) namedYarnCatalogs() ([]string, error) {
	path := filepath.Join(s.root, ".yarnrc.yml")
	payload, err := s.loadYAML(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	return sortedMapKeys(anyToStringMapMap(payload["npmCatalogs"])), nil
}

func (s Store) catalogPackageNamesBun(name string) ([]string, error) {
	path := filepath.Join(s.root, "package.json")
	content, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}

	payload := map[string]any{}
	if err := json.Unmarshal(content, &payload); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if strings.TrimSpace(name) == "" {
		return sortedStringMapKeys(anyToStringMap(payload["catalog"])), nil
	}
	catalogs := anyToStringMapMap(payload["catalogs"])
	return sortedStringMapKeys(catalogs[name]), nil
}

func (s Store) catalogPackageNamesPNPM(name string) ([]string, error) {
	path := filepath.Join(s.root, "pnpm-workspace.yaml")
	payload, err := s.loadYAML(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	if strings.TrimSpace(name) == "" {
		return sortedStringMapKeys(anyToStringMap(payload["catalog"])), nil
	}
	catalogs := anyToStringMapMap(payload["catalogs"])
	return sortedStringMapKeys(catalogs[name]), nil
}

func (s Store) catalogPackageNamesYarn(name string) ([]string, error) {
	path := filepath.Join(s.root, ".yarnrc.yml")
	payload, err := s.loadYAML(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	if strings.TrimSpace(name) == "" {
		return sortedStringMapKeys(anyToStringMap(payload["npmCatalog"])), nil
	}
	catalogs := anyToStringMapMap(payload["npmCatalogs"])
	return sortedStringMapKeys(catalogs[name]), nil
}

func (s Store) loadYAML(path string) (map[string]any, error) {
	content, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]any{}, nil
		}
		return nil, err
	}
	if len(content) == 0 {
		return map[string]any{}, nil
	}

	payload := map[string]any{}
	if err := yaml.Unmarshal(content, &payload); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if payload == nil {
		return map[string]any{}, nil
	}
	return payload, nil
}

func (s Store) writeYAML(path string, payload map[string]any) error {
	formatted, err := yaml.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	return s.fs.WriteFile(path, formatted, 0o644)
}

func upsertStringMap(current map[string]string, updates map[string]string, force bool) (map[string]string, error) {
	if current == nil {
		current = map[string]string{}
	}
	for name, version := range updates {
		existing, ok := current[name]
		if ok && existing != version && !force {
			return nil, fmt.Errorf("catalog conflict for %s: existing=%s want=%s", name, existing, version)
		}
		current[name] = version
	}
	return current, nil
}

func removeFromStringMap(current map[string]string, packages []string) map[string]string {
	if current == nil {
		return map[string]string{}
	}
	for _, pkg := range packages {
		delete(current, pkg)
	}
	return current
}

func anyToStringMap(value any) map[string]string {
	raw, ok := value.(map[string]any)
	if !ok {
		return map[string]string{}
	}
	out := make(map[string]string, len(raw))
	for key, value := range raw {
		str, ok := value.(string)
		if !ok {
			continue
		}
		out[key] = str
	}
	return out
}

func anyToStringMapMap(value any) map[string]map[string]string {
	raw, ok := value.(map[string]any)
	if !ok {
		return map[string]map[string]string{}
	}

	out := make(map[string]map[string]string, len(raw))
	for key, nested := range raw {
		out[key] = anyToStringMap(nested)
	}
	return out
}

func sortedMapKeys(items map[string]map[string]string) []string {
	keys := make([]string, 0, len(items))
	for key := range items {
		if strings.TrimSpace(key) == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedStringMapKeys(items map[string]string) []string {
	keys := make([]string, 0, len(items))
	for key := range items {
		if strings.TrimSpace(key) == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
