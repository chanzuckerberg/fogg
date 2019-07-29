package exp

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/exp/versioning"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	ExpCmd.AddCommand(versioningCmd)
}

var versioningCmd = &cobra.Command{
	Use:   "versioning",
	Short: "Detects terraform versioning changes",
	Long: `This command aims to detect changes between local terraform files
	and remote registries.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set up fs
		pwd, err := os.Getwd()
		if err != nil {
			return errs.WrapUser(err, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		return versioning.V(fs)
	},
}
