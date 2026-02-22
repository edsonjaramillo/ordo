package cli

import (
	"ordo/internal/app"
	"ordo/internal/cli/completion"
	"ordo/internal/cli/output"

	"github.com/spf13/cobra"
)

func newRunCmd(uc app.RunUseCase, completer completion.TargetCompleter, printer output.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <target> [-- <args...>]",
		Short: "Run a package script in root or workspace",
		Args:  cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			items, err := completer.ScriptTargets(cmd.Context(), toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return items, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uc.Run(cmd.Context(), app.RunRequest{Target: args[0], ExtraArgs: trailingArgs(args)})
			return printer.Handle(err)
		},
	}
	cmd.DisableFlagParsing = false
	return cmd
}

func trailingArgs(args []string) []string {
	if len(args) <= 1 {
		return nil
	}
	return args[1:]
}
