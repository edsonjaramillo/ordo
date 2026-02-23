package cli

import (
	"errors"
	"os"

	execadapter "ordo/internal/adapters/exec"
	fsadapter "ordo/internal/adapters/fs"
	registryadapter "ordo/internal/adapters/registry"
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
	suggestor := registryadapter.NewNPMSuggestor()
	printer := output.NewPrinter()
	configStore := fsadapter.NewConfigStore()
	installCompletion := app.NewInstallCompletionService(discovery, suggestor)
	completer := completion.NewTargetCompleter(discovery, installCompletion)
	globalCompletion := app.NewGlobalCompletionService(installCompletion, runner, runner)
	globalCompleter := completion.NewGlobalCompleter(globalCompletion)
	presetCompletion := app.NewPresetCompletionService(configStore)
	presetCompleter := completion.NewPresetCompleter(presetCompletion)

	runUC := app.NewRunUseCase(discovery, runner)
	installUC := app.NewInstallUseCase(discovery, runner)
	uninstallUC := app.NewUninstallUseCase(discovery, runner)
	updateUC := app.NewUpdateUseCase(discovery, runner)
	globalInstallUC := app.NewGlobalInstallUseCase(runner)
	globalUninstallUC := app.NewGlobalUninstallUseCase(runner, runner)
	globalUpdateUC := app.NewGlobalUpdateUseCase(runner)
	initUC := app.NewInitUseCase(configStore)
	presetUC := app.NewPresetUseCase(discovery, runner, configStore)
	var colorFlag string
	var noLevelFlag bool

	cmd := &cobra.Command{
		Use:           "ordo",
		Short:         "Run, install, uninstall, and update packages/scripts across JS monorepos",
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
	cmd.AddCommand(newInstallCmd(installUC, completer, printer))
	cmd.AddCommand(newUninstallCmd(uninstallUC, completer, printer))
	cmd.AddCommand(newUpdateCmd(updateUC, completer, printer))
	cmd.AddCommand(newGlobalCmd(globalInstallUC, globalUninstallUC, globalUpdateUC, globalCompleter, printer))
	cmd.AddCommand(newInitCmd(initUC, globalCompleter, printer))
	cmd.AddCommand(newPresetCmd(presetUC, presetCompleter, completer, printer))

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
