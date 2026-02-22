package domain

import "path/filepath"

type PackageInfo struct {
	Dir          string
	WorkspaceKey string
	Scripts      map[string]string
	Dependencies map[string]struct{}
	Lockfiles    map[string]bool
}

func WorkspaceKeyFromDir(dir string) string {
	if dir == "." || dir == "" {
		return ""
	}
	return filepath.Base(dir)
}
