package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"ordo/internal/cli/output"
)

func withRootOutputDefaults(t *testing.T) {
	t.Helper()
	output.SetOutputColorMode("auto")
	output.SetOutputShowLevel(true)
}

func newTestRootCmd(t *testing.T) (*cobra.Command, *bytes.Buffer) {
	t.Helper()
	withRootOutputDefaults(t)

	cmd, err := NewRootCmd()
	if err != nil {
		t.Fatalf("NewRootCmd() error = %v", err)
	}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	return cmd, buf
}

func TestRootRejectsInvalidColorFlagValue(t *testing.T) {
	cmd, _ := newTestRootCmd(t)
	cmd.SetArgs([]string{"--color=invalid", "run", "script"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), `invalid value for --color: "invalid"`) {
		t.Fatalf("error = %q, want invalid --color message", err.Error())
	}
}

func TestRootRejectsInvalidNoLevelEnvValue(t *testing.T) {
	t.Setenv("ORDO_NO_LEVEL", "maybe")
	cmd, _ := newTestRootCmd(t)
	cmd.SetArgs([]string{"run", "script"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), `invalid value for ORDO_NO_LEVEL: "maybe"`) {
		t.Fatalf("error = %q, want invalid ORDO_NO_LEVEL message", err.Error())
	}
}

func TestRootNoLevelFromEnvAppliesToSubcommands(t *testing.T) {
	t.Setenv("ORDO_NO_LEVEL", "true")
	cmd, buf := newTestRootCmd(t)
	cmd.SetArgs([]string{"run", "   "})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want non-nil")
	}
	if got := buf.String(); strings.Contains(got, "[ERROR]") {
		t.Fatalf("output = %q, want no level tag", got)
	}
}

func TestRootNoLevelFlagFalseOverridesEnvTrue(t *testing.T) {
	t.Setenv("ORDO_NO_LEVEL", "true")
	cmd, buf := newTestRootCmd(t)
	cmd.SetArgs([]string{"--no-level=false", "run", "   "})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want non-nil")
	}
	if got := buf.String(); !strings.Contains(got, "[ERROR]") {
		t.Fatalf("output = %q, want level tag", got)
	}
}

func TestRootNoLevelFlagTrueOverridesEnvFalse(t *testing.T) {
	t.Setenv("ORDO_NO_LEVEL", "false")
	cmd, buf := newTestRootCmd(t)
	cmd.SetArgs([]string{"--no-level=true", "run", "   "})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want non-nil")
	}
	if got := buf.String(); strings.Contains(got, "[ERROR]") {
		t.Fatalf("output = %q, want no level tag", got)
	}
}

func TestRootRegistersGlobalSubcommands(t *testing.T) {
	cmd, _ := newTestRootCmd(t)

	globalCmd, _, err := cmd.Find([]string{"global"})
	if err != nil {
		t.Fatalf("Find(global) error = %v", err)
	}
	if globalCmd == nil || globalCmd.Name() != "global" {
		t.Fatalf("expected global command, got %#v", globalCmd)
	}

	installCmd, _, err := cmd.Find([]string{"global", "install"})
	if err != nil {
		t.Fatalf("Find(global install) error = %v", err)
	}
	if installCmd == nil || installCmd.Name() != "install" {
		t.Fatalf("expected global install command, got %#v", installCmd)
	}

	uninstallCmd, _, err := cmd.Find([]string{"global", "uninstall"})
	if err != nil {
		t.Fatalf("Find(global uninstall) error = %v", err)
	}
	if uninstallCmd == nil || uninstallCmd.Name() != "uninstall" {
		t.Fatalf("expected global uninstall command, got %#v", uninstallCmd)
	}
}

func TestGlobalInstallRequiresManagerArg(t *testing.T) {
	cmd, _ := newTestRootCmd(t)
	cmd.SetArgs([]string{"global", "install", "typescript"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "at least 2 arg(s)") {
		t.Fatalf("error = %q, want arg count message", err.Error())
	}
}

func TestGlobalUninstallRejectsUnsupportedManager(t *testing.T) {
	cmd, buf := newTestRootCmd(t)
	cmd.SetArgs([]string{"global", "uninstall", "foobar", "typescript"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want non-nil")
	}
	if got := buf.String(); !strings.Contains(got, "unsupported package manager") {
		t.Fatalf("output = %q, want unsupported package manager message", got)
	}
}
