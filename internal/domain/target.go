package domain

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidTarget = errors.New("invalid target")

type Target struct {
	Workspace string
	Name      string
}

func ParseTarget(raw string) (Target, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return Target{}, fmt.Errorf("%w: target cannot be empty", ErrInvalidTarget)
	}

	parts := strings.Split(trimmed, "/")
	switch len(parts) {
	case 1:
		if parts[0] == "" {
			return Target{}, fmt.Errorf("%w: target cannot be empty", ErrInvalidTarget)
		}
		return Target{Name: parts[0]}, nil
	case 2:
		if parts[0] == "" || parts[1] == "" {
			return Target{}, fmt.Errorf("%w: workspace/name cannot contain empty segments", ErrInvalidTarget)
		}
		return Target{Workspace: parts[0], Name: parts[1]}, nil
	default:
		return Target{}, fmt.Errorf("%w: expected <name> or <workspace>/<name>", ErrInvalidTarget)
	}
}

func (t Target) IsRoot() bool {
	return t.Workspace == ""
}
