package cli

import (
	"ordo/internal/app"
	"ordo/internal/cli/completion"
	"ordo/internal/cli/output"

	"github.com/spf13/cobra"
)

func newUninstallCmd(uc app.UninstallUseCase, completer completion.TargetCompleter, printer output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall <target>",
		Short: "Uninstall a dependency in root or workspace",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			items, err := completer.PackageTargets(cmd.Context(), toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return items, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uc.Run(cmd.Context(), app.UninstallRequest{Target: args[0]})
			return printer.Handle(cmd.ErrOrStderr(), err)
		},
	}
}
