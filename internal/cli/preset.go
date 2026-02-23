package cli

import (
	"ordo/internal/app"
	"ordo/internal/cli/completion"
	"ordo/internal/cli/output"

	"github.com/spf13/cobra"
)

func newPresetCmd(
	uc app.PresetUseCase,
	completer completion.PresetCompleter,
	targets completion.TargetCompleter,
	printer output.Printer,
) *cobra.Command {
	var workspace string

	cmd := &cobra.Command{
		Use:   "preset <name> <bucket> [pkg[@version]...]",
		Short: "Install dependencies from a configured preset",
		Args:  cobra.MinimumNArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) {
			case 0:
				items, err := completer.PresetNames(cmd.Context(), toComplete)
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}
				return items, cobra.ShellCompDirectiveNoFileComp
			case 1:
				items, err := completer.Buckets(cmd.Context(), args[0], toComplete)
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}
				return items, cobra.ShellCompDirectiveNoFileComp
			default:
				items, err := completer.BucketPackages(cmd.Context(), args[0], args[1], toComplete)
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}
				return items, cobra.ShellCompDirectiveNoFileComp
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uc.Run(cmd.Context(), app.PresetRequest{
				Preset:    args[0],
				Bucket:    args[1],
				Packages:  args[2:],
				Workspace: workspace,
			})
			return printer.Handle(cmd.ErrOrStderr(), err)
		},
	}

	cmd.Flags().StringVar(&workspace, "workspace", "", "Workspace key to install into (default: root)")
	_ = cmd.RegisterFlagCompletionFunc("workspace", func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		items, err := targets.WorkspaceKeys(cmd.Context(), toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return items, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
