package domain

import (
	"fmt"
	"strings"
)

type PackageManager string

const (
	ManagerNPM  PackageManager = "npm"
	ManagerPNPM PackageManager = "pnpm"
	ManagerYarn PackageManager = "yarn"
	ManagerBun  PackageManager = "bun"
)

func SupportedPackageManagers() []string {
	return []string{
		string(ManagerNPM),
		string(ManagerPNPM),
		string(ManagerYarn),
		string(ManagerBun),
	}
}

func ParsePackageManager(raw string) (PackageManager, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	switch value {
	case string(ManagerNPM):
		return ManagerNPM, nil
	case string(ManagerPNPM):
		return ManagerPNPM, nil
	case string(ManagerYarn):
		return ManagerYarn, nil
	case string(ManagerBun):
		return ManagerBun, nil
	default:
		return "", fmt.Errorf("unsupported package manager: %q (supported: %s)", raw, strings.Join(SupportedPackageManagers(), ", "))
	}
}

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

func BuildUpdateCommand(manager PackageManager, pkg string) ([]string, error) {
	if pkg == "" {
		return nil, fmt.Errorf("package cannot be empty")
	}

	switch manager {
	case ManagerNPM:
		return []string{"npm", "update", pkg}, nil
	case ManagerPNPM:
		return []string{"pnpm", "update", pkg}, nil
	case ManagerYarn:
		return []string{"yarn", "upgrade", pkg}, nil
	case ManagerBun:
		return []string{"bun", "update", pkg}, nil
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", manager)
	}
}

type InstallOptions struct {
	Dev      bool
	Peer     bool
	Optional bool
	Prod     bool
	Exact    bool
}

func (o InstallOptions) validate() error {
	selected := 0
	for _, v := range []bool{o.Dev, o.Peer, o.Optional, o.Prod} {
		if v {
			selected++
		}
	}
	if selected > 1 {
		return fmt.Errorf("only one dependency type flag can be set")
	}
	return nil
}

func BuildInstallCommand(manager PackageManager, pkgs []string, opts InstallOptions) ([]string, error) {
	if err := validatePackages(pkgs); err != nil {
		return nil, err
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}

	var cmd []string
	switch manager {
	case ManagerNPM:
		cmd = []string{"npm", "install"}
		if opts.Dev {
			cmd = append(cmd, "--save-dev")
		}
		if opts.Peer {
			cmd = append(cmd, "--save-peer")
		}
		if opts.Optional {
			cmd = append(cmd, "--save-optional")
		}
		if opts.Prod {
			cmd = append(cmd, "--save-prod")
		}
		if opts.Exact {
			cmd = append(cmd, "--save-exact")
		}
	case ManagerPNPM:
		cmd = []string{"pnpm", "add"}
		if opts.Dev {
			cmd = append(cmd, "--save-dev")
		}
		if opts.Peer {
			cmd = append(cmd, "--save-peer")
		}
		if opts.Optional {
			cmd = append(cmd, "--save-optional")
		}
		if opts.Prod {
			cmd = append(cmd, "--save-prod")
		}
		if opts.Exact {
			cmd = append(cmd, "--save-exact")
		}
	case ManagerYarn:
		cmd = []string{"yarn", "add"}
		if opts.Dev {
			cmd = append(cmd, "--dev")
		}
		if opts.Peer {
			cmd = append(cmd, "--peer")
		}
		if opts.Optional {
			cmd = append(cmd, "--optional")
		}
		if opts.Prod {
			cmd = append(cmd, "--prod")
		}
		if opts.Exact {
			cmd = append(cmd, "--exact")
		}
	case ManagerBun:
		cmd = []string{"bun", "add"}
		if opts.Dev {
			cmd = append(cmd, "--dev")
		}
		if opts.Peer {
			cmd = append(cmd, "--peer")
		}
		if opts.Optional {
			cmd = append(cmd, "--optional")
		}
		if opts.Prod {
			cmd = append(cmd, "--production")
		}
		if opts.Exact {
			cmd = append(cmd, "--exact")
		}
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", manager)
	}
	return append(cmd, pkgs...), nil
}

func validatePackages(pkgs []string) error {
	if len(pkgs) == 0 {
		return fmt.Errorf("at least one package is required")
	}
	for _, pkg := range pkgs {
		if pkg == "" {
			return fmt.Errorf("package cannot be empty")
		}
	}
	return nil
}

func BuildGlobalInstallCommand(manager PackageManager, pkgs []string) ([]string, error) {
	if err := validatePackages(pkgs); err != nil {
		return nil, err
	}

	var cmd []string
	switch manager {
	case ManagerNPM:
		cmd = []string{"npm", "install", "--global"}
	case ManagerPNPM:
		cmd = []string{"pnpm", "add", "--global"}
	case ManagerYarn:
		cmd = []string{"yarn", "global", "add"}
	case ManagerBun:
		cmd = []string{"bun", "add", "--global"}
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", manager)
	}

	return append(cmd, pkgs...), nil
}

func BuildGlobalUninstallCommand(manager PackageManager, pkgs []string) ([]string, error) {
	if err := validatePackages(pkgs); err != nil {
		return nil, err
	}

	var cmd []string
	switch manager {
	case ManagerNPM:
		cmd = []string{"npm", "uninstall", "--global"}
	case ManagerPNPM:
		cmd = []string{"pnpm", "remove", "--global"}
	case ManagerYarn:
		cmd = []string{"yarn", "global", "remove"}
	case ManagerBun:
		cmd = []string{"bun", "remove", "--global"}
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", manager)
	}

	return append(cmd, pkgs...), nil
}

func BuildGlobalUpdateCommand(manager PackageManager, pkgs []string) ([]string, error) {
	if err := validatePackages(pkgs); err != nil {
		return nil, err
	}

	var cmd []string
	switch manager {
	case ManagerNPM:
		cmd = []string{"npm", "update", "--global"}
	case ManagerPNPM:
		cmd = []string{"pnpm", "update", "--global"}
	case ManagerYarn:
		cmd = []string{"yarn", "global", "upgrade"}
	case ManagerBun:
		cmd = []string{"bun", "update", "--global"}
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", manager)
	}

	return append(cmd, pkgs...), nil
}
