package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	planCmd.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	planCmd.Flags().BoolP("verbose", "v", false, "use this to turn on verbose output")
	rootCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run a plan",
	Long:  "plan will read fogg.json, use that to generate a plan and print that plan out. It will make no changes.",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel := log.InfoLevel
		if debug { // debug overrides quiet
			logLevel = log.DebugLevel
		} else if quiet {
			logLevel = log.FatalLevel
		}
		log.SetLevel(logLevel)

		var e error
		// Set up fs
		pwd, e := os.Getwd()
		if e != nil {
			log.Panic(e)
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		// handle flags
		verbose, e := cmd.Flags().GetBool("verbose")
		if e != nil {
			log.Panic(e)
		}

		configFile, e := cmd.Flags().GetString("config")
		if e != nil {
			log.Panic(e)
		}

		// check that we are at root of initialized git repo
		openGitOrExit(pwd)

		config, err := readAndValidateConfig(fs, configFile, verbose)

		exitOnConfigErrors(err)

		p, e := plan.Eval(config, verbose)
		if e != nil {
			log.Panic(e)
		}
		e = plan.Print(p)
		if e != nil {
			log.Panic(e)
		}
	},
}
