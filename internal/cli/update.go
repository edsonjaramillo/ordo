package cli

import (
	"ordo/internal/app"
	"ordo/internal/cli/completion"
	"ordo/internal/cli/output"

	"github.com/spf13/cobra"
)

func newUpdateCmd(uc app.UpdateUseCase, completer completion.TargetCompleter, printer output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "update <target>",
		Short: "Update a dependency in root or workspace",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			items, err := completer.PackageTargets(cmd.Context(), toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return items, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uc.Run(cmd.Context(), app.UpdateRequest{Target: args[0]})
			return printer.Handle(cmd.ErrOrStderr(), err)
		},
	}
}
