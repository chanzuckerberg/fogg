package cmd

import (
	"os"

	fogg_init "github.com/chanzuckerberg/fogg/init"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Run init",
	Run: func(cmd *cobra.Command, args []string) {
		var e error
		pwd, e := os.Getwd()
		if e != nil {
			panic(e)
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		e = fogg_init.PromptAndInit(fs)
		if e != nil {
			panic(e)
		}
	},
}
