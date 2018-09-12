package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
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
	RunE: func(cmd *cobra.Command, args []string) error {
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
			return errs.WrapUser(e, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		// handle flags
		verbose, e := cmd.Flags().GetBool("verbose")
		if e != nil {
			return errs.WrapInternal(e, "couldn't parse verbose flag")
		}

		configFile, e := cmd.Flags().GetString("config")
		if e != nil {
			return errs.WrapInternal(e, "couldn't parse config flag")
		}

		// check that we are at root of initialized git repo
		openGitOrExit(pwd)

		config, err := readAndValidateConfig(fs, configFile, verbose)

		exitOnConfigErrors(err)

		p, e := plan.Eval(config, verbose)
		if e != nil {
			return e
		}
		return plan.Print(p)
	},
}
