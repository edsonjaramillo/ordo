package output

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"ordo/internal/app"
	"ordo/internal/domain"
)

func withOutputColorMode(t *testing.T, mode string) {
	t.Helper()
	prev := outputColorMode
	SetOutputColorMode(mode)
	t.Cleanup(func() {
		SetOutputColorMode(prev)
	})
}

func withOutputShowLevel(t *testing.T, show bool) {
	t.Helper()
	prev := outputShowLevel
	SetOutputShowLevel(show)
	t.Cleanup(func() {
		SetOutputShowLevel(prev)
	})
}

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

func TestParseColorMode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "auto", input: "auto", want: colorModeAuto},
		{name: "always", input: "always", want: colorModeAlways},
		{name: "never", input: "never", want: colorModeNever},
		{name: "trim and case-insensitive", input: "  AlWaYs ", want: colorModeAlways},
		{name: "invalid", input: "sometimes", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseColorMode(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ParseColorMode(%q) error = nil, want non-nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseColorMode(%q) error = %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("ParseColorMode(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseNoLevelEnv(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		setEnv  bool
		wantNo  bool
		wantSet bool
		wantErr bool
	}{
		{name: "unset", setEnv: false, wantNo: false, wantSet: false, wantErr: false},
		{name: "true", value: "true", setEnv: true, wantNo: true, wantSet: true, wantErr: false},
		{name: "one", value: "1", setEnv: true, wantNo: true, wantSet: true, wantErr: false},
		{name: "false", value: "false", setEnv: true, wantNo: false, wantSet: true, wantErr: false},
		{name: "invalid", value: "maybe", setEnv: true, wantNo: false, wantSet: true, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			prev, hadPrev := os.LookupEnv(envNoLevel)
			t.Cleanup(func() {
				if hadPrev {
					_ = os.Setenv(envNoLevel, prev)
					return
				}
				_ = os.Unsetenv(envNoLevel)
			})

			if tc.setEnv {
				t.Setenv(envNoLevel, tc.value)
			} else {
				_ = os.Unsetenv(envNoLevel)
			}
			gotNoLevel, gotSet, err := ParseNoLevelEnv()
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ParseNoLevelEnv() error = nil, want non-nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseNoLevelEnv() error = %v", err)
			}
			if gotNoLevel != tc.wantNo {
				t.Fatalf("ParseNoLevelEnv() no-level = %v, want %v", gotNoLevel, tc.wantNo)
			}
			if gotSet != tc.wantSet {
				t.Fatalf("ParseNoLevelEnv() set = %v, want %v", gotSet, tc.wantSet)
			}
		})
	}
}

func TestPrintRootErrorWithLevelTag(t *testing.T) {
	withOutputColorMode(t, colorModeNever)
	withOutputShowLevel(t, true)

	stderr := &bytes.Buffer{}
	if err := PrintRootError(stderr, errors.New("boom")); err != nil {
		t.Fatalf("PrintRootError() error = %v", err)
	}
	if got, want := stderr.String(), "[ERROR] boom\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
}

func TestPrintRootErrorWithoutLevelTag(t *testing.T) {
	withOutputColorMode(t, colorModeAlways)
	withOutputShowLevel(t, false)

	stderr := &bytes.Buffer{}
	if err := PrintRootError(stderr, errors.New("boom")); err != nil {
		t.Fatalf("PrintRootError() error = %v", err)
	}
	if got, want := stderr.String(), "boom\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
}

func TestPrintRootErrorAlwaysIgnoresNoColorsEnv(t *testing.T) {
	withOutputColorMode(t, colorModeAlways)
	withOutputShowLevel(t, true)
	t.Setenv(envColorDisable, "1")

	stderr := &bytes.Buffer{}
	if err := PrintRootError(stderr, errors.New("boom")); err != nil {
		t.Fatalf("PrintRootError() error = %v", err)
	}
	if got, want := stderr.String(), "["+ansiRed+"ERROR"+ansiReset+"] boom\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
}

func TestPrintRootErrorAutoHonorsNoColorsEnv(t *testing.T) {
	withOutputColorMode(t, colorModeAuto)
	withOutputShowLevel(t, true)
	t.Setenv(envColorDisable, "1")

	stderr := &bytes.Buffer{}
	if err := PrintRootError(stderr, errors.New("boom")); err != nil {
		t.Fatalf("PrintRootError() error = %v", err)
	}
	if got, want := stderr.String(), "[ERROR] boom\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
}
