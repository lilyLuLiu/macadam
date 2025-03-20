package common

import "github.com/spf13/cobra"

var LogLevels = []string{"trace", "debug", "info", "warn", "warning", "error", "fatal", "panic"}

// AutocompleteLogLevel - Autocomplete log level options.
// -> "trace", "debug", "info", "warn", "error", "fatal", "panic"
func AutocompleteLogLevel(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return LogLevels, cobra.ShellCompDirectiveNoFileComp
}
