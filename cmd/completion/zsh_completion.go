package completion

import (
	"os"

	"github.com/spf13/cobra"
)

var zshCompletion = &cobra.Command{
	Use:   "zsh",
	Short: "Generates zsh completion",
	Long:  "Generates zsh completion",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.GenZshCompletion(os.Stdout)
	},
}

func init() {
	CompletionCmd.AddCommand(zshCompletion)
}
