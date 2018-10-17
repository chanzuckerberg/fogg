package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:          "completion",
	Short:        "",
	SilenceUsage: true,
}

var bashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generates bash completion",
	Long:  "Generates bash completion",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenBashCompletion(os.Stdout)
	},
}

var zshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generates zsh completion",
	Long:  "Generates zsh completion",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenZshCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(bashCmd)
	completionCmd.AddCommand(zshCmd)
}
