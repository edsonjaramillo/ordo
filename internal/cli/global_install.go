package cli

import (
	"ordo/internal/app"
	"ordo/internal/cli/completion"
	"ordo/internal/cli/output"
	"ordo/internal/domain"

	"github.com/spf13/cobra"
)

func newGlobalInstallCmd(uc app.GlobalInstallUseCase, completer completion.GlobalCompleter, printer output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "install <manager> <pkg[@version]>...",
		Short: "Install one or more global packages",
		Args:  cobra.MinimumNArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return domain.SupportedPackageManagers(), cobra.ShellCompDirectiveNoFileComp
			}
			if _, err := domain.ParsePackageManager(args[0]); err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			items, err := completer.InstallPackages(cmd.Context(), toComplete)
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

			err = uc.Run(cmd.Context(), app.GlobalInstallRequest{
				Manager:  manager,
				Packages: args[1:],
			})
			return printer.Handle(cmd.ErrOrStderr(), err)
		},
	}
}
