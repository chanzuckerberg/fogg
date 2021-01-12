package exp

import (
	"github.com/spf13/cobra"
)

func init() {
	ExpCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "",
	Long:  ``,
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	return nil
}
