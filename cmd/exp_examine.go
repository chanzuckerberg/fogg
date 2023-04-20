package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/exp/examine"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	expCmd.AddCommand(examineCmd)
}

// TODO:(EC) Create a flag for path to walk
var examineCmd = &cobra.Command{
	Use:   "examine",
	Short: "Detects terraform module updates",
	Long: `This command aims to detect version changes between local terraform file modules
	and remote registry modules.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set up fs
		pwd, err := os.Getwd()
		if err != nil {
			return errs.WrapUser(err, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		openGitOrExit(fs)

		return examine.Examine(fs, pwd)
	},
}
