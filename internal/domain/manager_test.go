package domain

import (
	"strings"
	"testing"
)

func TestDetectManagerPrecedence(t *testing.T) {
	tests := []struct {
		name      string
		lockfiles map[string]bool
		want      PackageManager
	}{
		{name: "bun wins", lockfiles: map[string]bool{"bun.lockb": true, "pnpm-lock.yaml": true}, want: ManagerBun},
		{name: "pnpm over yarn", lockfiles: map[string]bool{"pnpm-lock.yaml": true, "yarn.lock": true}, want: ManagerPNPM},
		{name: "yarn over npm", lockfiles: map[string]bool{"yarn.lock": true, "package-lock.json": true}, want: ManagerYarn},
		{name: "npm", lockfiles: map[string]bool{"package-lock.json": true}, want: ManagerNPM},
		{name: "default npm", lockfiles: map[string]bool{}, want: ManagerNPM},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := DetectManager(tc.lockfiles)
			if got != tc.want {
				t.Fatalf("expected %s, got %s", tc.want, got)
			}
		})
	}
}

func TestBuildInstallCommand(t *testing.T) {
	tests := []struct {
		name    string
		manager PackageManager
		pkgs    []string
		opts    InstallOptions
		want    []string
		wantErr bool
	}{
		{
			name:    "npm dev exact",
			manager: ManagerNPM,
			pkgs:    []string{"vite", "react@18"},
			opts:    InstallOptions{Dev: true, Exact: true},
			want:    []string{"npm", "install", "--save-dev", "--save-exact", "vite", "react@18"},
		},
		{
			name:    "pnpm peer",
			manager: ManagerPNPM,
			pkgs:    []string{"eslint"},
			opts:    InstallOptions{Peer: true},
			want:    []string{"pnpm", "add", "--save-peer", "eslint"},
		},
		{
			name:    "yarn optional exact",
			manager: ManagerYarn,
			pkgs:    []string{"left-pad"},
			opts:    InstallOptions{Optional: true, Exact: true},
			want:    []string{"yarn", "add", "--optional", "--exact", "left-pad"},
		},
		{
			name:    "bun prod",
			manager: ManagerBun,
			pkgs:    []string{"lodash"},
			opts:    InstallOptions{Prod: true},
			want:    []string{"bun", "add", "--production", "lodash"},
		},
		{
			name:    "conflicting type flags",
			manager: ManagerPNPM,
			pkgs:    []string{"pkg"},
			opts:    InstallOptions{Dev: true, Optional: true},
			wantErr: true,
		},
		{
			name:    "empty package list",
			manager: ManagerPNPM,
			pkgs:    nil,
			opts:    InstallOptions{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := BuildInstallCommand(tc.manager, tc.pkgs, tc.opts)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("BuildInstallCommand() error = nil, want non-nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("BuildInstallCommand() error = %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("BuildInstallCommand() len = %d, want %d (%#v)", len(got), len(tc.want), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("BuildInstallCommand()[%d] = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestBuildGlobalInstallCommand(t *testing.T) {
	tests := []struct {
		name    string
		manager PackageManager
		pkgs    []string
		want    []string
		wantErr bool
	}{
		{
			name:    "npm",
			manager: ManagerNPM,
			pkgs:    []string{"typescript"},
			want:    []string{"npm", "install", "--global", "typescript"},
		},
		{
			name:    "pnpm",
			manager: ManagerPNPM,
			pkgs:    []string{"typescript"},
			want:    []string{"pnpm", "add", "--global", "typescript"},
		},
		{
			name:    "yarn",
			manager: ManagerYarn,
			pkgs:    []string{"typescript"},
			want:    []string{"yarn", "global", "add", "typescript"},
		},
		{
			name:    "bun",
			manager: ManagerBun,
			pkgs:    []string{"typescript"},
			want:    []string{"bun", "add", "--global", "typescript"},
		},
		{
			name:    "empty package list",
			manager: ManagerNPM,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := BuildGlobalInstallCommand(tc.manager, tc.pkgs)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("BuildGlobalInstallCommand() error = nil, want non-nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("BuildGlobalInstallCommand() error = %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("BuildGlobalInstallCommand() len = %d, want %d (%#v)", len(got), len(tc.want), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("BuildGlobalInstallCommand()[%d] = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestBuildGlobalUninstallCommand(t *testing.T) {
	tests := []struct {
		name    string
		manager PackageManager
		pkgs    []string
		want    []string
		wantErr bool
	}{
		{
			name:    "npm",
			manager: ManagerNPM,
			pkgs:    []string{"typescript"},
			want:    []string{"npm", "uninstall", "--global", "typescript"},
		},
		{
			name:    "pnpm",
			manager: ManagerPNPM,
			pkgs:    []string{"typescript"},
			want:    []string{"pnpm", "remove", "--global", "typescript"},
		},
		{
			name:    "yarn",
			manager: ManagerYarn,
			pkgs:    []string{"typescript"},
			want:    []string{"yarn", "global", "remove", "typescript"},
		},
		{
			name:    "bun",
			manager: ManagerBun,
			pkgs:    []string{"typescript"},
			want:    []string{"bun", "remove", "--global", "typescript"},
		},
		{
			name:    "empty package list",
			manager: ManagerNPM,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := BuildGlobalUninstallCommand(tc.manager, tc.pkgs)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("BuildGlobalUninstallCommand() error = nil, want non-nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("BuildGlobalUninstallCommand() error = %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("BuildGlobalUninstallCommand() len = %d, want %d (%#v)", len(got), len(tc.want), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("BuildGlobalUninstallCommand()[%d] = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestBuildUpdateCommand(t *testing.T) {
	tests := []struct {
		name    string
		manager PackageManager
		pkg     string
		want    []string
		wantErr bool
	}{
		{
			name:    "npm",
			manager: ManagerNPM,
			pkg:     "typescript",
			want:    []string{"npm", "update", "typescript"},
		},
		{
			name:    "pnpm",
			manager: ManagerPNPM,
			pkg:     "typescript",
			want:    []string{"pnpm", "update", "typescript"},
		},
		{
			name:    "yarn",
			manager: ManagerYarn,
			pkg:     "typescript",
			want:    []string{"yarn", "upgrade", "typescript"},
		},
		{
			name:    "bun",
			manager: ManagerBun,
			pkg:     "typescript",
			want:    []string{"bun", "update", "typescript"},
		},
		{
			name:    "empty package",
			manager: ManagerNPM,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := BuildUpdateCommand(tc.manager, tc.pkg)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("BuildUpdateCommand() error = nil, want non-nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("BuildUpdateCommand() error = %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("BuildUpdateCommand() len = %d, want %d (%#v)", len(got), len(tc.want), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("BuildUpdateCommand()[%d] = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestBuildGlobalUpdateCommand(t *testing.T) {
	tests := []struct {
		name    string
		manager PackageManager
		pkgs    []string
		want    []string
		wantErr bool
	}{
		{
			name:    "npm",
			manager: ManagerNPM,
			pkgs:    []string{"typescript"},
			want:    []string{"npm", "update", "--global", "typescript"},
		},
		{
			name:    "pnpm",
			manager: ManagerPNPM,
			pkgs:    []string{"typescript"},
			want:    []string{"pnpm", "update", "--global", "typescript"},
		},
		{
			name:    "yarn",
			manager: ManagerYarn,
			pkgs:    []string{"typescript"},
			want:    []string{"yarn", "global", "upgrade", "typescript"},
		},
		{
			name:    "bun",
			manager: ManagerBun,
			pkgs:    []string{"typescript"},
			want:    []string{"bun", "update", "--global", "typescript"},
		},
		{
			name:    "empty package list",
			manager: ManagerNPM,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := BuildGlobalUpdateCommand(tc.manager, tc.pkgs)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("BuildGlobalUpdateCommand() error = nil, want non-nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("BuildGlobalUpdateCommand() error = %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("BuildGlobalUpdateCommand() len = %d, want %d (%#v)", len(got), len(tc.want), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("BuildGlobalUpdateCommand()[%d] = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestSupportedPackageManagers(t *testing.T) {
	got := SupportedPackageManagers()
	want := []string{"npm", "pnpm", "yarn", "bun"}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestParsePackageManager(t *testing.T) {
	tests := []struct {
		in      string
		want    PackageManager
		wantErr bool
	}{
		{in: "npm", want: ManagerNPM},
		{in: "PNPM", want: ManagerPNPM},
		{in: " yarn ", want: ManagerYarn},
		{in: "bun", want: ManagerBun},
		{in: "invalid", wantErr: true},
	}

	for _, tc := range tests {
		got, err := ParsePackageManager(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("ParsePackageManager(%q) expected error", tc.in)
			}
			if !strings.Contains(err.Error(), "supported: npm, pnpm, yarn, bun") {
				t.Fatalf("unexpected error: %v", err)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ParsePackageManager(%q) error: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("ParsePackageManager(%q) = %s, want %s", tc.in, got, tc.want)
		}
	}
}
