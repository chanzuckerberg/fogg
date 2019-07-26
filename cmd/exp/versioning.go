package exp

import (
	"github.com/chanzuckerberg/fogg/exp/versioning"
	"github.com/spf13/cobra"
)

func init() {
	ExpCmd.AddCommand(versioning)
}

var migrateCmd = &cobra.Command{
	Use:   "versioning",
	Short: "Detects terraform versioning changes",
	Long: `This command aims to detect changes between local terraform files
	and remote registries.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//TODO: Return an actual value
		return nil
	)
	},
}
