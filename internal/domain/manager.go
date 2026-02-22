package domain

import "fmt"

type PackageManager string

const (
	ManagerNPM  PackageManager = "npm"
	ManagerPNPM PackageManager = "pnpm"
	ManagerYarn PackageManager = "yarn"
	ManagerBun  PackageManager = "bun"
)

func DetectManager(lockfiles map[string]bool) PackageManager {
	if lockfiles["bun.lockb"] || lockfiles["bun.lock"] {
		return ManagerBun
	}
	if lockfiles["pnpm-lock.yaml"] {
		return ManagerPNPM
	}
	if lockfiles["yarn.lock"] {
		return ManagerYarn
	}
	if lockfiles["package-lock.json"] || lockfiles["npm-shrinkwrap.json"] {
		return ManagerNPM
	}
	return ManagerNPM
}

func BuildRunCommand(manager PackageManager, script string, extraArgs []string) ([]string, error) {
	if script == "" {
		return nil, fmt.Errorf("script cannot be empty")
	}

	var cmd []string
	switch manager {
	case ManagerNPM:
		cmd = []string{"npm", "run", script}
	case ManagerPNPM:
		cmd = []string{"pnpm", "run", script}
	case ManagerYarn:
		cmd = []string{"yarn", "run", script}
	case ManagerBun:
		cmd = []string{"bun", "run", script}
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", manager)
	}
	return append(cmd, extraArgs...), nil
}

func BuildUninstallCommand(manager PackageManager, pkg string) ([]string, error) {
	if pkg == "" {
		return nil, fmt.Errorf("package cannot be empty")
	}

	switch manager {
	case ManagerNPM:
		return []string{"npm", "uninstall", pkg}, nil
	case ManagerPNPM:
		return []string{"pnpm", "remove", pkg}, nil
	case ManagerYarn:
		return []string{"yarn", "remove", pkg}, nil
	case ManagerBun:
		return []string{"bun", "remove", pkg}, nil
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", manager)
	}
}
