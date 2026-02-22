package output

import (
	"errors"
	"testing"

	"ordo/internal/app"
	"ordo/internal/domain"
)

func TestExitCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "invalid target", err: domain.ErrInvalidTarget, want: 2},
		{name: "workspace missing", err: app.ErrWorkspaceNotFound, want: 3},
		{name: "script missing", err: app.ErrScriptNotFound, want: 4},
		{name: "package missing", err: app.ErrPackageNotFound, want: 4},
		{name: "other", err: errors.New("boom"), want: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ExitCode(tc.err); got != tc.want {
				t.Fatalf("expected %d, got %d", tc.want, got)
			}
		})
	}
}
