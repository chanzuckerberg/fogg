package exp

import "github.com/spf13/cobra"

// ExpCmd is a subcommand for experimental commands
var ExpCmd = &cobra.Command{
	Use:          "exp",
	Short:        "Experimental commands",
	Long:         "Grouping of experimental commands. These are experimental and prone to change",
	SilenceUsage: true,
}
