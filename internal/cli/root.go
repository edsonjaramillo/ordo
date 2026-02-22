package cli

import (
	"errors"
	"os"

	execadapter "ordo/internal/adapters/exec"
	fsadapter "ordo/internal/adapters/fs"
	"ordo/internal/app"
	"ordo/internal/cli/completion"
	"ordo/internal/cli/output"
	"ordo/internal/config"

	"github.com/spf13/cobra"
)

func NewRootCmd() (*cobra.Command, error) {
	cwd, err := config.WorkingDir()
	if err != nil {
		return nil, err
	}

	indexer := fsadapter.NewWorkspaceIndexer(cwd)
	discovery := app.NewDiscoveryService(indexer)
	runner := execadapter.NewRunner()
	printer := output.NewPrinter(os.Stderr)
	completer := completion.NewTargetCompleter(discovery)

	runUC := app.NewRunUseCase(discovery, runner)
	uninstallUC := app.NewUninstallUseCase(discovery, runner)

	cmd := &cobra.Command{
		Use:          "ordo",
		Short:        "Run and uninstall packages/scripts across JS monorepos",
		SilenceUsage: true,
	}

	cmd.AddCommand(newRunCmd(runUC, completer, printer))
	cmd.AddCommand(newUninstallCmd(uninstallUC, completer, printer))

	return cmd, nil
}

func Execute() int {
	cmd, err := NewRootCmd()
	if err != nil {
		_, _ = os.Stderr.WriteString("Error: " + err.Error() + "\n")
		return 1
	}

	err = cmd.Execute()
	if err == nil {
		return 0
	}

	var exitErr *output.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.Code
	}
	_, _ = os.Stderr.WriteString("Error: " + err.Error() + "\n")
	return 1
}
