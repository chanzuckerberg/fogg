package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/apply"
	"github.com/chanzuckerberg/fogg/templates"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	applyCmd.Flags().BoolP("sicc", "s", false, "Use this to turn on sicc-compatibility mode. Implies -c sicc.json.")
	applyCmd.Flags().BoolP("verbose", "v", false, "use this to turn on verbose output")
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
			log.Panic(e)
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		// handle flags
		siccMode, e := cmd.Flags().GetBool("sicc")
		if e != nil {
			log.Panic(e)
		}
		verbose, e := cmd.Flags().GetBool("verbose")
		if e != nil {
			log.Panic(e)
		}
		var configFile string
		if siccMode {
			configFile = "sicc.json"
		} else {
			configFile, e = cmd.Flags().GetString("config")
			if e != nil {
				log.Panic(e)
			}
		}

		// check that we are at root of initialized git repo
		openGitOrExit(pwd)

		config, err := readAndValidateConfig(fs, configFile, verbose)

		exitOnConfigErrors(err)

		// apply
		e = apply.Apply(fs, config, templates.Templates, siccMode)
		if e != nil {
			log.Panic(e)
		}
	},
}
