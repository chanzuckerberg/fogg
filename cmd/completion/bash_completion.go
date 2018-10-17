package completion

import (
	"os"

	"github.com/spf13/cobra"
)

var bashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generates bash completion",
	Long:  "Generates bash completion",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	CompletionCmd.AddCommand(bashCmd)
}
