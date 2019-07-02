package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/migrations"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	migrateCommand.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	migrateCommand.Flags().BoolP("skip", "s", false, "Use this to run all tests.")
	rootCmd.AddCommand(migrateCommand)
}

var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Converts existing fogg.json to fogg.yml",
	Long:  "This command will convert the fogg.json to a fogg.yml file type.",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Set up fs
		pwd, err := os.Getwd()
		if err != nil {
			return errs.WrapUser(err, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			return errs.WrapInternal(err, "couldn't parse config flag")
		}

		skipPrompts, err := cmd.Flags().GetBool("skip")
		if err != nil {
			return errs.WrapInternal(err, "couldn't parse skip flag")
		}

		openGitOrExit(fs)

		return migrations.RunMigrations(fs, configFile, skipPrompts)
	},
}
