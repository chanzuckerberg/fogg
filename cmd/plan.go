package cmd

import (
	"os"

	"github.com/ryanking/fogg/plan"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run a plan",
	Run: func(cmd *cobra.Command, args []string) {
		pwd, _ := os.Getwd()
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		p, _ := plan.Plan(fs)
		plan.Print(p)
	},
}
