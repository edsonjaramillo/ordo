package cli

import (
	"ordo/internal/app"
	"ordo/internal/cli/completion"
	"ordo/internal/cli/output"

	"github.com/spf13/cobra"
)

func newGlobalCmd(
	installUC app.GlobalInstallUseCase,
	uninstallUC app.GlobalUninstallUseCase,
	updateUC app.GlobalUpdateUseCase,
	completer completion.GlobalCompleter,
	printer output.Printer,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "global",
		Short: "Manage global packages",
	}

	cmd.AddCommand(newGlobalInstallCmd(installUC, completer, printer))
	cmd.AddCommand(newGlobalUninstallCmd(uninstallUC, completer, printer))
	cmd.AddCommand(newGlobalUpdateCmd(updateUC, completer, printer))

	return cmd
}
