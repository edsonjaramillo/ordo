package cli

import (
	"ordo/internal/app"
	"ordo/internal/cli/completion"
	"ordo/internal/cli/output"

	"github.com/spf13/cobra"
)

func newInitCmd(uc app.InitUseCase, completer completion.GlobalCompleter, printer output.Printer) *cobra.Command {
	var defaultPackageManager string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create ordo config in XDG config home",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := uc.Run(cmd.Context(), app.InitRequest{DefaultPackageManager: defaultPackageManager})
			return printer.Handle(cmd.ErrOrStderr(), err)
		},
	}

	cmd.Flags().StringVar(&defaultPackageManager, "defaultPackageManager", "", "Default package manager for generated config (bun, npm, pnpm, yarn)")
	_ = cmd.MarkFlagRequired("defaultPackageManager")
	_ = cmd.RegisterFlagCompletionFunc("defaultPackageManager", func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		items, err := completer.AvailablePackageManagers(cmd.Context(), toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return items, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
