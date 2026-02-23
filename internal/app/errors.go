package app

import (
	"errors"
	"fmt"
	"strings"

	"ordo/internal/domain"
)

var (
	ErrWorkspaceNotFound   = errors.New("workspace not found")
	ErrScriptNotFound      = errors.New("script not found")
	ErrPackageNotFound     = errors.New("package not found")
	ErrConfigAlreadyExists = errors.New("ordo config already exists")
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
