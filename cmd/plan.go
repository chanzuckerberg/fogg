package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/plan"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	planCmd.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	planCmd.Flags().BoolP("sicc", "s", false, "Use this to turn on sicc-compatibility mode. Implies -c sicc.json.")
	planCmd.Flags().BoolP("verbose", "v", false, "use this to turn on verbose output")
	rootCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run a plan",
	Long:  "plan will read fogg.json, use that to generate a plan and print that plan out. It will make no changes.",
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
		verbose, e := cmd.Flags().GetBool("verbose")
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

		config, err := readAndValidateConfig(fs, configFile, verbose)

		exitOnConfigErrors(err)

		p, e := plan.Eval(config, siccMode, verbose)
		if e != nil {
			panic(e)
		}
		e = plan.Print(p)
		if e != nil {
			panic(e)
		}
	},
}
