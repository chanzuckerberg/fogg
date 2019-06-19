package exp

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/exp/yaml_migrate"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	ExpCmd.AddCommand(yamlMigrateCmd)
}

var yamlMigrateCmd = &cobra.Command{
	Use:   "yaml-migrate",
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

		yaml_migrate.OpenGitOrExit(fs)

		return yaml_migrate.ConvertToYaml(fs, configFile)
	},
}
