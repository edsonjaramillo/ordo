package cli

import (
	"ordo/internal/app"
	"ordo/internal/cli/completion"
	"ordo/internal/cli/output"

	"github.com/spf13/cobra"
)

func newCatalogCmd(
	uc app.CatalogUseCase,
	catalogCompleter completion.CatalogCompleter,
	presetCompleter completion.PresetCompleter,
	printer output.Printer,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "Manage default catalog entries",
	}

	cmd.AddCommand(newCatalogAddCmd(uc, catalogCompleter, printer))
	cmd.AddCommand(newCatalogImportCmd(uc, catalogCompleter, printer))
	cmd.AddCommand(newCatalogPresetsCmd(uc, catalogCompleter, presetCompleter, printer))
	cmd.AddCommand(newCatalogRemoveCmd(uc, catalogCompleter, printer))
	cmd.AddCommand(newCatalogSyncCmd(uc, printer))

	return cmd
}

func newCatalogImportCmd(uc app.CatalogUseCase, completer completion.CatalogCompleter, printer output.Printer) *cobra.Command {
	var fromWorkspace string
	var force bool

	cmd := &cobra.Command{
		Use:   "import <pkg>",
		Short: "Import a workspace dependency version into the root catalog",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			items, err := completer.WorkspaceDependencyNames(cmd.Context(), fromWorkspace, toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return items, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uc.RunImport(cmd.Context(), app.CatalogImportRequest{
				Package:       args[0],
				FromWorkspace: fromWorkspace,
				Force:         force,
			})
			return printer.Handle(cmd.ErrOrStderr(), err)
		},
	}

	cmd.Flags().StringVar(&fromWorkspace, "from-workspace", "", "Source workspace key to copy package version from")
	cmd.Flags().BoolVar(&force, "force", false, "Override conflicting existing catalog versions")
	mustMarkFlagRequired(cmd, "from-workspace")
	mustRegisterFlagCompletionFunc(cmd, "from-workspace", func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		items, err := completer.WorkspaceKeys(cmd.Context(), toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return items, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func newCatalogAddCmd(uc app.CatalogUseCase, completer completion.CatalogCompleter, printer output.Printer) *cobra.Command {
	var workspace string
	var force bool

	cmd := &cobra.Command{
		Use:   "add <pkg[@version]>...",
		Short: "Add default catalog entries and rewrite dependency references",
		Args:  cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			items, err := completer.PackageSpecs(cmd.Context(), toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			items = filterCompletedArgs(items, args, 0)
			return items, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uc.RunAdd(cmd.Context(), app.CatalogAddRequest{
				Packages:  args,
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

func newCatalogRemoveCmd(uc app.CatalogUseCase, completer completion.CatalogCompleter, printer output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <pkg>...",
		Short: "Remove package entries from the default catalog",
		Args:  cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			items, err := completer.CatalogPackageNames(cmd.Context(), "", toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			items = filterCompletedArgs(items, args, 0)
			return items, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uc.RunRemove(cmd.Context(), app.CatalogRemoveRequest{Packages: args})
			return printer.Handle(cmd.ErrOrStderr(), err)
		},
	}
}

func newCatalogSyncCmd(uc app.CatalogUseCase, printer output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sync root catalog references to all nested workspaces",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := uc.RunSync(cmd.Context(), app.CatalogSyncRequest{})
			return printer.Handle(cmd.ErrOrStderr(), err)
		},
	}
}

func newCatalogPresetsCmd(
	uc app.CatalogUseCase,
	catalogCompleter completion.CatalogCompleter,
	presetCompleter completion.PresetCompleter,
	printer output.Printer,
) *cobra.Command {
	var workspace string
	var force bool

	cmd := &cobra.Command{
		Use:   "presets <name> <bucket> [pkg[@version]...]",
		Short: "Add preset packages to the default catalog and rewrite dependency references",
		Args:  cobra.MinimumNArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) {
			case 0:
				items, err := presetCompleter.PresetNames(cmd.Context(), toComplete)
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}
				return items, cobra.ShellCompDirectiveNoFileComp
			case 1:
				items, err := presetCompleter.Buckets(cmd.Context(), args[0], toComplete)
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}
				return items, cobra.ShellCompDirectiveNoFileComp
			default:
				items, err := presetCompleter.BucketPackages(cmd.Context(), args[0], args[1], toComplete)
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}
				items = filterCompletedArgs(items, args, 2)
				return items, cobra.ShellCompDirectiveNoFileComp
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uc.RunPreset(cmd.Context(), app.CatalogPresetRequest{
				Preset:    args[0],
				Bucket:    args[1],
				Packages:  args[2:],
				Workspace: workspace,
				Force:     force,
			})
			return printer.Handle(cmd.ErrOrStderr(), err)
		},
	}

	cmd.Flags().StringVar(&workspace, "workspace", "", "Workspace key to apply dependency references (default: root)")
	cmd.Flags().BoolVar(&force, "force", false, "Override conflicting existing catalog versions")
	mustRegisterFlagCompletionFunc(cmd, "workspace", func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		items, err := catalogCompleter.WorkspaceKeys(cmd.Context(), toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return items, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
