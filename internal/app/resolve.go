package app

import (
	"fmt"
	"strings"

	"ordo/internal/domain"
)

func resolveTargetPackage(snapshot Snapshot, target domain.Target) (domain.PackageInfo, error) {
	if target.IsRoot() {
		return snapshot.Root, nil
	}

	pkg, ok := snapshot.ByWorkspace[target.Workspace]
	if !ok {
		return domain.PackageInfo{}, fmt.Errorf("%w: %s", ErrWorkspaceNotFound, target.Workspace)
	}
	return pkg, nil
}

func resolveInstallTargetPackage(snapshot Snapshot, workspace string) (domain.PackageInfo, error) {
	key := strings.TrimSpace(workspace)
	if key == "" {
		return snapshot.Root, nil
	}
	pkg, ok := snapshot.ByWorkspace[key]
	if !ok {
		return domain.PackageInfo{}, fmt.Errorf("%w: %s", ErrWorkspaceNotFound, key)
	}
	return pkg, nil
}
