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
	applyCmd.Flags().BoolP("sicc", "s", false, "Use this to turn on sicc-compatibility mode. Implies -c sicc.json.")
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run an apply",
	Run: func(cmd *cobra.Command, args []string) {
		var e error
		// Set up fs
		pwd, e := os.Getwd()
		if e != nil {
			panic(e)
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		// handle flags
		siccMode, e := cmd.Flags().GetBool("sicc")
		if e != nil {
			panic(e)
		}
		var configFile string
		if siccMode {
			configFile = "sicc.json"
		} else {
			configFile, e = cmd.Flags().GetString("config")
			if e != nil {
				panic(e)
			}
		}

		// check that we are at root of initialized git repo
		openGitOrExit(pwd)

		// apply
		e = apply.Apply(fs, configFile, templates.Templates, siccMode)
		if e != nil {
			panic(e)
		}
	},
}
