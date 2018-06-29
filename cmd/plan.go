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
	rootCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run a plan",
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

		p, e := plan.Eval(fs, configFile, siccMode)
		if e != nil {
			panic(e)
		}
		e = plan.Print(p)
		if e != nil {
			panic(e)
		}
	},
}
