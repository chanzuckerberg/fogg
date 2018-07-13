package cmd

import (
	"os"

	fogg_init "github.com/chanzuckerberg/fogg/init"
	log "github.com/sirupsen/logrus"
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
			log.Panic(e)
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		// check that we are at root of initialized git repo
		openGitOrExit(pwd)

		e = fogg_init.Init(fs)
		if e != nil {
			log.Panic(e)
		}
	},
}
