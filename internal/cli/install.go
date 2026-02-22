package cli

import (
	"ordo/internal/app"
	"ordo/internal/cli/completion"
	"ordo/internal/cli/output"

	"github.com/spf13/cobra"
)

func newInstallCmd(uc app.InstallUseCase, completer completion.TargetCompleter, printer output.Printer) *cobra.Command {
	var workspace string
	var dev bool
	var peer bool
	var optional bool
	var prod bool
	var exact bool

	cmd := &cobra.Command{
		Use:   "install <pkg[@version]>...",
		Short: "Install one or more dependencies in root or workspace",
		Args:  cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			items, err := completer.InstallPackages(cmd.Context(), toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return items, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uc.Run(cmd.Context(), app.InstallRequest{
				Packages:  args,
				Workspace: workspace,
				Dev:       dev,
				Peer:      peer,
				Optional:  optional,
				Prod:      prod,
				Exact:     exact,
			})
			return printer.Handle(cmd.ErrOrStderr(), err)
		},
	}

	cmd.Flags().StringVar(&workspace, "workspace", "", "Workspace key to install into (default: root)")
	cmd.Flags().BoolVar(&dev, "dev", false, "Install as a development dependency")
	cmd.Flags().BoolVar(&peer, "peer", false, "Install as a peer dependency")
	cmd.Flags().BoolVar(&optional, "optional", false, "Install as an optional dependency")
	cmd.Flags().BoolVar(&prod, "prod", false, "Install as a production dependency")
	cmd.Flags().BoolVar(&exact, "exact", false, "Pin exact version")

	_ = cmd.RegisterFlagCompletionFunc("workspace", func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		items, err := completer.WorkspaceKeys(cmd.Context(), toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return items, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
