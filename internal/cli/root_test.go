package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
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

	updateCmd, _, err := cmd.Find([]string{"global", "update"})
	if err != nil {
		t.Fatalf("Find(global update) error = %v", err)
	}
	if updateCmd == nil || updateCmd.Name() != "update" {
		t.Fatalf("expected global update command, got %#v", updateCmd)
	}
}

func TestRootRegistersUpdateCommand(t *testing.T) {
	cmd, _ := newTestRootCmd(t)

	updateCmd, _, err := cmd.Find([]string{"update"})
	if err != nil {
		t.Fatalf("Find(update) error = %v", err)
	}
	if updateCmd == nil || updateCmd.Name() != "update" {
		t.Fatalf("expected update command, got %#v", updateCmd)
	}
}

func TestRootRegistersInitCommand(t *testing.T) {
	cmd, _ := newTestRootCmd(t)

	initCmd, _, err := cmd.Find([]string{"init"})
	if err != nil {
		t.Fatalf("Find(init) error = %v", err)
	}
	if initCmd == nil || initCmd.Name() != "init" {
		t.Fatalf("expected init command, got %#v", initCmd)
	}
}

func TestRootRegistersPresetCommand(t *testing.T) {
	cmd, _ := newTestRootCmd(t)

	presetCmd, _, err := cmd.Find([]string{"preset"})
	if err != nil {
		t.Fatalf("Find(preset) error = %v", err)
	}
	if presetCmd == nil || presetCmd.Name() != "preset" {
		t.Fatalf("expected preset command, got %#v", presetCmd)
	}
}

func TestInitRequiresDefaultPackageManagerFlag(t *testing.T) {
	cmd, _ := newTestRootCmd(t)
	cmd.SetArgs([]string{"init"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), `required flag(s) "defaultPackageManager" not set`) {
		t.Fatalf("error = %q, want required defaultPackageManager flag message", err.Error())
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

func TestGlobalUpdateRequiresManagerArg(t *testing.T) {
	cmd, _ := newTestRootCmd(t)
	cmd.SetArgs([]string{"global", "update", "typescript"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "at least 2 arg(s)") {
		t.Fatalf("error = %q, want arg count message", err.Error())
	}
}

func TestGlobalInstallCompletionUsesPathAvailableManagers(t *testing.T) {
	pathDir := t.TempDir()
	writeTestExecutable(t, pathDir, "npm")
	t.Setenv("PATH", pathDir)

	cmd, _ := newTestRootCmd(t)
	installCmd, _, err := cmd.Find([]string{"global", "install"})
	if err != nil {
		t.Fatalf("Find(global install) error = %v", err)
	}

	items, dir := installCmd.ValidArgsFunction(installCmd, []string{}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(items) != 1 || items[0] != "npm" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestGlobalInstallCompletionFallsBackWhenNoManagersOnPath(t *testing.T) {
	pathDir := t.TempDir()
	t.Setenv("PATH", pathDir)

	cmd, _ := newTestRootCmd(t)
	installCmd, _, err := cmd.Find([]string{"global", "install"})
	if err != nil {
		t.Fatalf("Find(global install) error = %v", err)
	}

	items, _ := installCmd.ValidArgsFunction(installCmd, []string{}, "")
	want := []string{"bun", "npm", "pnpm", "yarn"}
	if len(items) != len(want) {
		t.Fatalf("unexpected items len: %#v", items)
	}
	for i := range want {
		if items[i] != want[i] {
			t.Fatalf("items[%d] = %q, want %q", i, items[i], want[i])
		}
	}
}

func TestInitDefaultPackageManagerCompletionUsesPathAvailableManagers(t *testing.T) {
	pathDir := t.TempDir()
	writeTestExecutable(t, pathDir, "pnpm")
	writeTestExecutable(t, pathDir, "yarn")
	t.Setenv("PATH", pathDir)

	cmd, _ := newTestRootCmd(t)
	initCmd, _, err := cmd.Find([]string{"init"})
	if err != nil {
		t.Fatalf("Find(init) error = %v", err)
	}

	fn, ok := initCmd.GetFlagCompletionFunc("defaultPackageManager")
	if !ok {
		t.Fatal("expected completion function for defaultPackageManager")
	}

	items, dir := fn(initCmd, []string{}, "y")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(items) != 1 || items[0] != "yarn" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func writeTestExecutable(t *testing.T, dir string, name string) {
	t.Helper()

	filename := name
	content := "#!/bin/sh\nexit 0\n"
	mode := os.FileMode(0o755)
	if runtime.GOOS == "windows" {
		filename += ".exe"
		content = ""
	}

	if err := os.WriteFile(filepath.Join(dir, filename), []byte(content), mode); err != nil {
		t.Fatalf("write executable %q: %v", filename, err)
	}
}
