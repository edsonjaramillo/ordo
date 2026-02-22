package fs

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"ordo/internal/domain"
)

var ignoredDirs = map[string]struct{}{
	".git":         {},
	"node_modules": {},
	"dist":         {},
	"build":        {},
	"coverage":     {},
	".next":        {},
}

type WorkspaceIndexer struct {
	root string
}

func NewWorkspaceIndexer(root string) WorkspaceIndexer {
	return WorkspaceIndexer{root: root}
}

func (w WorkspaceIndexer) Discover(ctx context.Context) ([]domain.PackageInfo, error) {
	infos := make([]domain.PackageInfo, 0)

	err := filepath.WalkDir(w.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if d.IsDir() {
			if _, skip := ignoredDirs[d.Name()]; skip {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Name() != "package.json" {
			return nil
		}

		dir := filepath.Dir(path)
		relDir, err := filepath.Rel(w.root, dir)
		if err != nil {
			return err
		}
		relDir = filepath.ToSlash(relDir)
		if relDir == "" {
			relDir = "."
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		pkg, err := parsePackageInfo(content)
		if err != nil {
			return err
		}
		pkg.Dir = relDir
		pkg.Lockfiles = map[string]bool{}

		if relDir == "." {
			for _, lockfile := range []string{"bun.lockb", "bun.lock", "pnpm-lock.yaml", "yarn.lock", "package-lock.json", "npm-shrinkwrap.json"} {
				if fileExists(filepath.Join(w.root, lockfile)) {
					pkg.Lockfiles[lockfile] = true
				}
			}
		}

		infos = append(infos, pkg)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return infos, nil
}

type packageJSON struct {
	Scripts              map[string]string `json:"scripts"`
	Dependencies         map[string]string `json:"dependencies"`
	DevDependencies      map[string]string `json:"devDependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
}

func parsePackageInfo(content []byte) (domain.PackageInfo, error) {
	var manifest packageJSON
	if err := json.Unmarshal(content, &manifest); err != nil {
		return domain.PackageInfo{}, err
	}

	deps := map[string]struct{}{}
	for name := range manifest.Dependencies {
		deps[name] = struct{}{}
	}
	for name := range manifest.DevDependencies {
		deps[name] = struct{}{}
	}
	for name := range manifest.PeerDependencies {
		deps[name] = struct{}{}
	}
	for name := range manifest.OptionalDependencies {
		deps[name] = struct{}{}
	}

	scripts := manifest.Scripts
	if scripts == nil {
		scripts = map[string]string{}
	}

	return domain.PackageInfo{
		Scripts:      scripts,
		Dependencies: deps,
	}, nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
