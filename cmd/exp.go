package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(expCmd)
}

// ExpCmd is a subcommand for experimental commands
var expCmd = &cobra.Command{
	Use:          "exp",
	Short:        "Experimental commands",
	Long:         "Grouping of experimental commands. These are experimental and prone to change",
	SilenceUsage: true,
}
