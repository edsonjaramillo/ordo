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
	printer := output.NewPrinter()
	completer := completion.NewTargetCompleter(discovery)

	runUC := app.NewRunUseCase(discovery, runner)
	uninstallUC := app.NewUninstallUseCase(discovery, runner)
	var colorFlag string
	var noLevelFlag bool

	cmd := &cobra.Command{
		Use:           "ordo",
		Short:         "Run and uninstall packages/scripts across JS monorepos",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			mode, err := output.ParseColorMode(colorFlag)
			if err != nil {
				return err
			}

			showLevel := true
			noLevelFromEnv, envSet, err := output.ParseNoLevelEnv()
			if err != nil {
				return err
			}
			if envSet {
				showLevel = !noLevelFromEnv
			}
			if cmd.Flags().Changed("no-level") {
				showLevel = !noLevelFlag
			}

			output.SetOutputColorMode(mode)
			output.SetOutputShowLevel(showLevel)
			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&colorFlag, "color", "auto", "Colorize output: auto, always, never")
	cmd.PersistentFlags().BoolVar(&noLevelFlag, "no-level", false, "Hide output level labels (INFO, OK, WARN, ERROR)")

	cmd.AddCommand(newRunCmd(runUC, completer, printer))
	cmd.AddCommand(newUninstallCmd(uninstallUC, completer, printer))

	return cmd, nil
}

func Execute() int {
	cmd, err := NewRootCmd()
	if err != nil {
		_ = output.PrintRootError(os.Stderr, err)
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
	_ = output.PrintRootError(os.Stderr, err)
	return 1
}
