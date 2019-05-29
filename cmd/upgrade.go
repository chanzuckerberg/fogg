package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	upgrade.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	rootCmd.AddCommand(upgrade)
}

var upgrade = &cobra.Command{
	Use:   "v2upgrade",
	Short: "Upgrades a v1 config to a v2 config",
	Long: `This command will upgrade a v1 config to a v2 config.
	Note that this is a lossy transformation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set up fs
		pwd, err := os.Getwd()
		if err != nil {
			return errs.WrapUser(err, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		// handle flags
		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			return errs.WrapInternal(err, "couldn't parse config flag")
		}

		// check that we are at root of initialized git repo
		openGitOrExit(fs)
		return config.Upgrade(fs, configFile)
	},
}
