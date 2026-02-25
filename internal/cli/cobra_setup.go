package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func mustRegisterFlagCompletionFunc(
	cmd *cobra.Command,
	name string,
	fn func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective),
) {
	if err := cmd.RegisterFlagCompletionFunc(name, fn); err != nil {
		panic(fmt.Sprintf("register completion for %s --%s: %v", cmd.CommandPath(), name, err))
	}
}

func mustMarkFlagRequired(cmd *cobra.Command, name string) {
	if err := cmd.MarkFlagRequired(name); err != nil {
		panic(fmt.Sprintf("mark required flag for %s --%s: %v", cmd.CommandPath(), name, err))
	}
}
