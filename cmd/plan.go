package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/plan"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var ConfigFile string

func init() {
	planCmd.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	rootCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run a plan",
	Run: func(cmd *cobra.Command, args []string) {
		pwd, _ := os.Getwd()
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
		configFile, _ := cmd.Flags().GetString("config")

		p, err := plan.Eval(fs, configFile)
		if err != nil {
			util.Dump(err)
			return
		}
		plan.Print(p)
	},
}
