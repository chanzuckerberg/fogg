package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/exp/state"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	autoRemoteStateCmd.Flags().String("path", "", "path to a working directory")
	autoRemoteStateCmd.Flags().StringP("config", "c", "fogg.yml", "Use this to override the fogg config file.")

	expCmd.AddCommand(autoRemoteStateCmd)
}

var autoRemoteStateCmd = &cobra.Command{
	Use:   "auto-remote-state",
	Short: "Read all the code in a given directory and update fogg.yml with depends_on configuration.",
	Long: `Read all the code in a given directory and update fogg.yml with depends_on configuration.

	Example usage- fogg exp auto-remote-state --path terraform/envs/prod/snowalert

	BEWARE– This is a very experimental feature, not well tested and with rough edges.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set up fs
		pwd, err := os.Getwd()
		if err != nil {
			return errs.WrapUser(err, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		configFile, e := cmd.Flags().GetString("config")
		if e != nil {
			return errs.WrapInternal(e, "couldn't parse config flag")
		}

		openGitOrExit(fs)

		path, _ := cmd.Flags().GetString("path")
		return state.Run(fs, configFile, path)
	},
}
