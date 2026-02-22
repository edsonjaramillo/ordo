package app

import "errors"

var (
	ErrWorkspaceNotFound = errors.New("workspace not found")
	ErrScriptNotFound    = errors.New("script not found")
	ErrPackageNotFound   = errors.New("package not found")
)
