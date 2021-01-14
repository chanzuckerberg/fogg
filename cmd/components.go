package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(componentsCmd)
}

// ComponentsCmd is a subcommand working with components
var componentsCmd = &cobra.Command{
	Use:          "components",
	SilenceUsage: true,
}
