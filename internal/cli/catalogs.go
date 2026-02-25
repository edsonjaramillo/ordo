package cli

import (
	"ordo/internal/app"
	"ordo/internal/cli/completion"
	"ordo/internal/cli/output"

	"github.com/spf13/cobra"
)

func newCatalogsCmd(uc app.CatalogUseCase, completer completion.CatalogCompleter, printer output.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "catalogs",
		Short: "Manage named catalogs",
	}

	cmd.AddCommand(newCatalogsAddCmd(uc, completer, printer))
	cmd.AddCommand(newCatalogsRemoveCmd(uc, completer, printer))

	return cmd
}

func newCatalogsAddCmd(uc app.CatalogUseCase, completer completion.CatalogCompleter, printer output.Printer) *cobra.Command {
	var workspace string
	var force bool

	cmd := &cobra.Command{
		Use:   "add <name> <pkg[@version]>...",
		Short: "Add named catalog entries and rewrite dependency references",
		Args:  cobra.MinimumNArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				items, err := completer.NamedCatalogs(cmd.Context(), toComplete)
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}
				return items, cobra.ShellCompDirectiveNoFileComp
			}

			items, err := completer.PackageSpecs(cmd.Context(), toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			items = filterCompletedArgs(items, args, 1)
			return items, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uc.RunNamedAdd(cmd.Context(), app.CatalogsAddRequest{
				Name:      args[0],
				Packages:  args[1:],
				Workspace: workspace,
				Force:     force,
			})
			return printer.Handle(cmd.ErrOrStderr(), err)
		},
	}

	cmd.Flags().StringVar(&workspace, "workspace", "", "Workspace key to apply dependency references (default: root)")
	cmd.Flags().BoolVar(&force, "force", false, "Override conflicting existing catalog versions")
	mustRegisterFlagCompletionFunc(cmd, "workspace", func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		items, err := completer.WorkspaceKeys(cmd.Context(), toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return items, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func newCatalogsRemoveCmd(uc app.CatalogUseCase, completer completion.CatalogCompleter, printer output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name> <pkg>...",
		Short: "Remove package entries from a named catalog",
		Args:  cobra.MinimumNArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				items, err := completer.NamedCatalogs(cmd.Context(), toComplete)
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}
				return items, cobra.ShellCompDirectiveNoFileComp
			}

			items, err := completer.CatalogPackageNames(cmd.Context(), args[0], toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			items = filterCompletedArgs(items, args, 1)
			return items, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uc.RunNamedRemove(cmd.Context(), app.CatalogsRemoveRequest{
				Name:     args[0],
				Packages: args[1:],
			})
			return printer.Handle(cmd.ErrOrStderr(), err)
		},
	}
}
