package exp

import (
	"github.com/chanzuckerberg/fogg/exp/yaml_migrate"
	"github.com/spf13/cobra"
)

func init() {
	ExpCmd.AddCommand(yamlMigrateCmd)
}

var yamlMigrateCmd = &cobra.Command{
	Use:   "yaml_migrate",
	Short: "Converts fogg.json to fogg.yml",
	Long:  "This command will convert the fogg.json to a fogg.yml file type.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return yaml_migrate.JSONtoYML()
	},
}
