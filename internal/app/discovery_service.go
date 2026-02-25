package app

import (
	"context"
	"path/filepath"
	"sort"
	"strings"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type DiscoveryService struct {
	indexer ports.WorkspaceIndexer
}

func NewDiscoveryService(indexer ports.WorkspaceIndexer) DiscoveryService {
	return DiscoveryService{indexer: indexer}
}

func (d DiscoveryService) Snapshot(ctx context.Context) (Snapshot, error) {
	infos, err := d.indexer.Discover(ctx)
	if err != nil {
		return Snapshot{}, err
	}

	prepared := make([]domain.PackageInfo, 0, len(infos))
	for _, info := range infos {
		if info.Scripts == nil {
			info.Scripts = map[string]string{}
		}
		if info.Dependencies == nil {
			info.Dependencies = map[string]struct{}{}
		}
		if info.DependencyVersions == nil {
			info.DependencyVersions = map[string]string{}
		}
		if info.Lockfiles == nil {
			info.Lockfiles = map[string]bool{}
		}
		prepared = append(prepared, info)
	}

	assignWorkspaceKeys(prepared)

	var root domain.PackageInfo
	byWorkspace := map[string]domain.PackageInfo{}
	lockfiles := map[string]bool{}

	for _, item := range prepared {
		if item.Dir == "." || item.Dir == "" {
			root = item
			continue
		}
		byWorkspace[item.WorkspaceKey] = item
	}

	for _, name := range []string{"bun.lockb", "bun.lock", "pnpm-lock.yaml", "yarn.lock", "package-lock.json", "npm-shrinkwrap.json"} {
		if root.Lockfiles[name] {
			lockfiles[name] = true
		}
	}

	return Snapshot{
		Root:        root,
		ByWorkspace: byWorkspace,
		Manager:     domain.DetectManager(lockfiles),
	}, nil
}

type Snapshot struct {
	Root        domain.PackageInfo
	ByWorkspace map[string]domain.PackageInfo
	Manager     domain.PackageManager
}

func assignWorkspaceKeys(infos []domain.PackageInfo) {
	counts := map[string]int{}
	for _, info := range infos {
		if info.Dir == "." || info.Dir == "" {
			continue
		}
		base := domain.WorkspaceKeyFromDir(info.Dir)
		counts[base]++
	}

	for i := range infos {
		if infos[i].Dir == "." || infos[i].Dir == "" {
			infos[i].WorkspaceKey = ""
			continue
		}
		base := domain.WorkspaceKeyFromDir(infos[i].Dir)
		if counts[base] > 1 {
			infos[i].WorkspaceKey = filepath.ToSlash(infos[i].Dir)
			continue
		}
		infos[i].WorkspaceKey = base
	}
}

func (s Snapshot) ScriptTargets(prefix string) []string {
	items := make([]string, 0)
	for name := range s.Root.Scripts {
		items = append(items, name)
	}
	for workspace, pkg := range s.ByWorkspace {
		for script := range pkg.Scripts {
			items = append(items, workspace+"/"+script)
		}
	}
	return filterAndSort(items, prefix)
}

func (s Snapshot) PackageTargets(prefix string) []string {
	items := make([]string, 0)
	for name := range s.Root.Dependencies {
		items = append(items, name)
	}
	for workspace, pkg := range s.ByWorkspace {
		for dep := range pkg.Dependencies {
			items = append(items, workspace+"/"+dep)
		}
	}
	return filterAndSort(items, prefix)
}

func (s Snapshot) WorkspaceKeys(prefix string) []string {
	items := make([]string, 0, len(s.ByWorkspace))
	for key := range s.ByWorkspace {
		items = append(items, key)
	}
	return filterAndSort(items, prefix)
}

func (s Snapshot) DependencyNames(prefix string) []string {
	items := make([]string, 0)
	for name := range s.Root.Dependencies {
		items = append(items, name)
	}
	for _, pkg := range s.ByWorkspace {
		for dep := range pkg.Dependencies {
			items = append(items, dep)
		}
	}
	return filterAndSort(items, prefix)
}

func filterAndSort(items []string, prefix string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if prefix != "" && !strings.HasPrefix(item, prefix) {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	sort.Strings(out)
	return out
}
