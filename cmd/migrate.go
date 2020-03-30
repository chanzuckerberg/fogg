package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/migrations"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	migrateCommand.Flags().StringP("config", "c", "fogg.yml", "Use this to override the fogg config file.")
	migrateCommand.Flags().BoolP("force", "f", false, "Use this to skip all of the migration prompts.")
	rootCmd.AddCommand(migrateCommand)
}

var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Runs all possible fogg migrations",
	Long:  "This command will run all applicable updates to fogg configuration.",
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

		skipPrompts, err := cmd.Flags().GetBool("force")
		if err != nil {
			return errs.WrapInternal(err, "couldn't parse skip flag")
		}

		openGitOrExit(fs)

		return migrations.RunMigrations(fs, configFile, skipPrompts)
	},
}
