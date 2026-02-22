package cli

import (
	"ordo/internal/app"
	"ordo/internal/cli/completion"
	"ordo/internal/cli/output"
	"ordo/internal/domain"

	"github.com/spf13/cobra"
)

func newGlobalUninstallCmd(uc app.GlobalUninstallUseCase, completer completion.GlobalCompleter, printer output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall <manager> <pkg>...",
		Short: "Uninstall one or more global packages",
		Args:  cobra.MinimumNArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return domain.SupportedPackageManagers(), cobra.ShellCompDirectiveNoFileComp
			}
			manager, err := domain.ParsePackageManager(args[0])
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			items, err := completer.InstalledGlobalPackages(cmd.Context(), manager, toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return items, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := domain.ParsePackageManager(args[0])
			if err != nil {
				return printer.Handle(cmd.ErrOrStderr(), err)
			}

			err = uc.Run(cmd.Context(), app.GlobalUninstallRequest{
				Manager:  manager,
				Packages: args[1:],
			})
			return printer.Handle(cmd.ErrOrStderr(), err)
		},
	}
}
