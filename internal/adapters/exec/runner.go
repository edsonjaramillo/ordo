package execadapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"ordo/internal/domain"
)

type Runner struct{}

func NewRunner() Runner {
	return Runner{}
}

func (r Runner) Run(ctx context.Context, dir string, argv []string) error {
	if len(argv) == 0 {
		return fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (r Runner) AvailablePackageManagers(_ context.Context) ([]string, error) {
	found := make([]string, 0, len(domain.SupportedPackageManagers()))
	for _, manager := range domain.SupportedPackageManagers() {
		if _, err := exec.LookPath(manager); err == nil {
			found = append(found, manager)
		}
	}
	sort.Strings(found)
	return found, nil
}

func (r Runner) ListInstalledGlobalPackages(ctx context.Context, manager domain.PackageManager) ([]string, error) {
	paths, err := r.ResolveGlobalStorePaths(ctx, manager)
	if err != nil {
		return nil, err
	}

	seen := map[string]struct{}{}
	for _, store := range paths {
		items, err := listPackagesFromNodeModules(store)
		if err != nil {
			continue
		}
		for _, item := range items {
			seen[item] = struct{}{}
		}
	}
	if len(seen) > 0 {
		return sortKeys(seen), nil
	}

	// Fall back to manager-native listing in case the global store structure differs.
	return listWithPackageManager(ctx, manager)
}

func (r Runner) ResolveGlobalStorePaths(ctx context.Context, manager domain.PackageManager) ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, 6)
	switch manager {
	case domain.ManagerNPM:
		if root, err := runCommandOutput(ctx, []string{"npm", "root", "-g"}); err == nil {
			paths = append(paths, root)
		}
		if prefix, err := runCommandOutput(ctx, []string{"npm", "config", "get", "prefix"}); err == nil {
			paths = append(paths, filepath.Join(prefix, "lib", "node_modules"))
		}
	case domain.ManagerPNPM:
		if root, err := runCommandOutput(ctx, []string{"pnpm", "root", "-g"}); err == nil {
			paths = append(paths, root)
		}
		if pnpmHome := strings.TrimSpace(os.Getenv("PNPM_HOME")); pnpmHome != "" {
			paths = append(paths, filepath.Join(pnpmHome, "global", "node_modules"))
			paths = append(paths, globNodeModulePaths(filepath.Join(pnpmHome, "global", "*", "node_modules"))...)
		}
		paths = append(paths, filepath.Join(home, ".local", "share", "pnpm", "global", "node_modules"))
		paths = append(paths, globNodeModulePaths(filepath.Join(home, ".local", "share", "pnpm", "global", "*", "node_modules"))...)
	case domain.ManagerYarn:
		if globalDir, err := runCommandOutput(ctx, []string{"yarn", "global", "dir"}); err == nil {
			paths = append(paths, filepath.Join(globalDir, "node_modules"))
		}
		if yarnGlobal := strings.TrimSpace(os.Getenv("YARN_GLOBAL_FOLDER")); yarnGlobal != "" {
			paths = append(paths, filepath.Join(yarnGlobal, "node_modules"))
		}
		paths = append(paths, filepath.Join(home, ".config", "yarn", "global", "node_modules"))
	case domain.ManagerBun:
		bunInstall := strings.TrimSpace(os.Getenv("BUN_INSTALL"))
		if bunInstall == "" {
			bunInstall = filepath.Join(home, ".bun")
		}
		paths = append(paths, filepath.Join(bunInstall, "install", "global", "node_modules"))
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", manager)
	}

	return uniqueNonEmptyPaths(paths), nil
}

func listWithPackageManager(ctx context.Context, manager domain.PackageManager) ([]string, error) {
	argv, err := listGlobalCommand(manager)
	if err != nil {
		return nil, err
	}
	stdout, err := runCommandOutput(ctx, argv)
	if err != nil {
		return nil, err
	}
	return parseListOutput(manager, stdout)
}

func runCommandOutput(ctx context.Context, argv []string) (string, error) {
	if len(argv) == 0 {
		return "", fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = "command failed"
		}
		return "", fmt.Errorf("%s: %w", msg, err)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func globNodeModulePaths(pattern string) []string {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}
	return matches
}

func uniqueNonEmptyPaths(paths []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		clean := filepath.Clean(p)
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		out = append(out, clean)
	}
	return out
}

func listPackagesFromNodeModules(nodeModulesDir string) ([]string, error) {
	entries, err := os.ReadDir(nodeModulesDir)
	if err != nil {
		return nil, err
	}

	seen := map[string]struct{}{}
	for _, entry := range entries {
		name := entry.Name()
		if name == "" || strings.HasPrefix(name, ".") {
			continue
		}

		if strings.HasPrefix(name, "@") {
			if !entry.IsDir() {
				continue
			}
			scopedEntries, err := os.ReadDir(filepath.Join(nodeModulesDir, name))
			if err != nil {
				continue
			}
			for _, scopedEntry := range scopedEntries {
				if !scopedEntry.IsDir() {
					continue
				}
				scopedName := strings.TrimSpace(scopedEntry.Name())
				if scopedName == "" || strings.HasPrefix(scopedName, ".") {
					continue
				}
				seen[name+"/"+scopedName] = struct{}{}
			}
			continue
		}

		if entry.IsDir() {
			seen[name] = struct{}{}
		}
	}

	return sortKeys(seen), nil
}

func listGlobalCommand(manager domain.PackageManager) ([]string, error) {
	switch manager {
	case domain.ManagerNPM:
		return []string{"npm", "ls", "--global", "--depth=0", "--json"}, nil
	case domain.ManagerPNPM:
		return []string{"pnpm", "ls", "--global", "--depth=0", "--json"}, nil
	case domain.ManagerYarn:
		return []string{"yarn", "global", "list", "--json"}, nil
	case domain.ManagerBun:
		return []string{"bun", "pm", "ls", "--global", "--json"}, nil
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", manager)
	}
}

func parseListOutput(manager domain.PackageManager, raw string) ([]string, error) {
	switch manager {
	case domain.ManagerYarn:
		return parseYarnGlobalListOutput(raw), nil
	case domain.ManagerNPM, domain.ManagerPNPM, domain.ManagerBun:
		return parseJSONDependenciesOutput(raw)
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", manager)
	}
}

func parseJSONDependenciesOutput(raw string) ([]string, error) {
	var parsed any
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil, err
	}

	seen := map[string]struct{}{}
	collectDependencyKeys(parsed, seen)
	return sortKeys(seen), nil
}

func collectDependencyKeys(v any, out map[string]struct{}) {
	switch val := v.(type) {
	case map[string]any:
		if depsRaw, ok := val["dependencies"]; ok {
			if deps, ok := depsRaw.(map[string]any); ok {
				for name := range deps {
					out[name] = struct{}{}
				}
			}
		}
		for _, child := range val {
			collectDependencyKeys(child, out)
		}
	case []any:
		for _, child := range val {
			collectDependencyKeys(child, out)
		}
	}
}

func parseYarnGlobalListOutput(raw string) []string {
	seen := map[string]struct{}{}
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}
		var payload map[string]any
		if err := json.Unmarshal([]byte(line), &payload); err != nil {
			continue
		}
		collectYarnNames(payload, seen)
	}
	return sortKeys(seen)
}

func collectYarnNames(v any, out map[string]struct{}) {
	switch val := v.(type) {
	case map[string]any:
		for key, child := range val {
			if key == "name" {
				if str, ok := child.(string); ok {
					name := trimVersionFromSpec(str)
					if name != "" {
						out[name] = struct{}{}
					}
				}
			}
			collectYarnNames(child, out)
		}
	case []any:
		for _, child := range val {
			collectYarnNames(child, out)
		}
	}
}

func trimVersionFromSpec(spec string) string {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return ""
	}
	i := strings.LastIndex(spec, "@")
	if i <= 0 {
		return spec
	}
	return spec[:i]
}

func sortKeys(in map[string]struct{}) []string {
	out := make([]string, 0, len(in))
	for item := range in {
		out = append(out, item)
	}
	sort.Strings(out)
	return out
}
