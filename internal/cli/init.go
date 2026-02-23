package cli

import (
	"ordo/internal/app"
	"ordo/internal/cli/output"
	"ordo/internal/domain"

	"github.com/spf13/cobra"
)

func newInitCmd(uc app.InitUseCase, printer output.Printer) *cobra.Command {
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
	_ = cmd.RegisterFlagCompletionFunc("defaultPackageManager", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return domain.SupportedPackageManagers(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
