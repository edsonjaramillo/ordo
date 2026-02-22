package domain

import "testing"

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
