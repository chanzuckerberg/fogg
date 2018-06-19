package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/apply"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run an apply",
	Run: func(cmd *cobra.Command, args []string) {
		pwd, _ := os.Getwd()
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
		configFile, _ := cmd.Flags().GetString("config")

		err := apply.Apply(fs, configFile, templates.Templates)
		if err != nil {
			util.Dump(err)
			return
		}
		// apply.Print(p)
	},
}
