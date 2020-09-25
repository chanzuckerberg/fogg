package exp

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/exp/state"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	autoRemoteStateCmd.Flags().String("path", "", "path to a working directory")
	ExpCmd.AddCommand(autoRemoteStateCmd)
}

//TODO:(EC) Create a flag for path to walk
var autoRemoteStateCmd = &cobra.Command{
	Use:   "auto-remote-state",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set up fs
		pwd, err := os.Getwd()
		if err != nil {
			return errs.WrapUser(err, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		openGitOrExit(fs)

		path, _ := cmd.Flags().GetString("path")
		return state.Run(fs, path)
	},
}
