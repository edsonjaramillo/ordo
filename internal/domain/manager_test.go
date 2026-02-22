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
