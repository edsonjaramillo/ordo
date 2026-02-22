package app

import (
	"fmt"

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
