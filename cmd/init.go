package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	fogg_init "github.com/chanzuckerberg/fogg/init"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new repo for use with fogg",
	Long:  "fogg init will ask you some questions and generate a basic fogg.yml.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var e error
		pwd, e := os.Getwd()
		if e != nil {
			return errs.WrapUser(e, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
		// check that we are at root of initialized git repo
		openGitOrExit(fs)

		return fogg_init.Init(fs)
	},
}
