package completion

import (
	"github.com/spf13/cobra"
)

// CompletionCmd represents the completion command
var CompletionCmd = &cobra.Command{
	Use:          "completion",
	Short:        "",
	SilenceUsage: true,
}
