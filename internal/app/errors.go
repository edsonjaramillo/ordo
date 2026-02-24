package app

import (
	"errors"
	"fmt"
	"strings"

	"ordo/internal/domain"
)

var (
	ErrWorkspaceNotFound     = errors.New("workspace not found")
	ErrScriptNotFound        = errors.New("script not found")
	ErrPackageNotFound       = errors.New("package not found")
	ErrConfigAlreadyExists   = errors.New("ordo config already exists")
	ErrConfigNotFound        = errors.New("ordo config not found")
	ErrPresetNotFound        = errors.New("preset not found")
	ErrPresetBucketNotFound  = errors.New("preset bucket not found")
	ErrPresetPackageNotFound = errors.New("preset package not found")
	ErrCatalogUnsupported    = errors.New("catalogs are unsupported for package manager")
	ErrCatalogConflict       = errors.New("catalog entry conflict")
	ErrInvalidCatalogName    = errors.New("invalid catalog name")
)

type GlobalPackageMissingError struct {
	Manager      domain.PackageManager
	Missing      []string
	CheckedPaths []string
}

func (e GlobalPackageMissingError) Error() string {
	msg := fmt.Sprintf("global package(s) not found for %s: %s", e.Manager, strings.Join(e.Missing, ", "))
	if len(e.CheckedPaths) == 0 {
		return msg
	}
	return msg + " (checked: " + strings.Join(e.CheckedPaths, ", ") + ")"
}
