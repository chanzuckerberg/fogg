package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/apply"
	"github.com/chanzuckerberg/fogg/templates"
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
		var e error
		pwd, e := os.Getwd()
		if e != nil {
			panic(e)
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
		configFile, e := cmd.Flags().GetString("config")
		if e != nil {
			panic(e)
		}

		e = apply.Apply(fs, configFile, templates.Templates)
		if e != nil {
			panic(e)
		}
		// apply.Print(p)
	},
}
