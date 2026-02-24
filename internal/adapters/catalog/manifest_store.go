package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type ManifestStore struct {
	root string
	fs   ports.ConfigStore
}

func NewManifestStore(root string, fs ports.ConfigStore) ManifestStore {
	return ManifestStore{root: root, fs: fs}
}

func (s ManifestStore) RewriteCatalogReferences(_ context.Context, targetDir string, catalogName string, packages []string) error {
	return s.rewriteCatalogReferences(targetDir, catalogName, packages, true)
}

func (s ManifestStore) RewriteCatalogReferencesExistingOnly(_ context.Context, targetDir string, catalogName string, packages []string) error {
	return s.rewriteCatalogReferences(targetDir, catalogName, packages, false)
}

func (s ManifestStore) rewriteCatalogReferences(targetDir string, catalogName string, packages []string, addMissing bool) error {
	path := filepath.Join(s.root, targetDir, "package.json")
	content, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("package manifest not found: %s", path)
		}
		return err
	}

	manifest := map[string]any{}
	if err := json.Unmarshal(content, &manifest); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	dependencies := anyToManifestMap(manifest["dependencies"])
	devDependencies := anyToManifestMap(manifest["devDependencies"])
	peerDependencies := anyToManifestMap(manifest["peerDependencies"])
	optionalDependencies := anyToManifestMap(manifest["optionalDependencies"])

	ref := domain.CatalogReference(catalogName)
	for _, pkg := range packages {
		updated := false
		if rewriteDependency(dependencies, pkg, ref) {
			updated = true
		}
		if rewriteDependency(devDependencies, pkg, ref) {
			updated = true
		}
		if rewriteDependency(peerDependencies, pkg, ref) {
			updated = true
		}
		if rewriteDependency(optionalDependencies, pkg, ref) {
			updated = true
		}
		if addMissing && !updated {
			dependencies[pkg] = ref
		}
	}

	manifest["dependencies"] = dependencies
	if len(devDependencies) > 0 {
		manifest["devDependencies"] = devDependencies
	}
	if len(peerDependencies) > 0 {
		manifest["peerDependencies"] = peerDependencies
	}
	if len(optionalDependencies) > 0 {
		manifest["optionalDependencies"] = optionalDependencies
	}

	formatted, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	formatted = append(formatted, '\n')
	return s.fs.WriteFile(path, formatted, 0o644)
}

func rewriteDependency(deps map[string]string, pkg string, ref string) bool {
	if _, ok := deps[pkg]; !ok {
		return false
	}
	deps[pkg] = ref
	return true
}

func anyToManifestMap(value any) map[string]string {
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
