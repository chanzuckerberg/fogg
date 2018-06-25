package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of fogg",
	Run: func(cmd *cobra.Command, args []string) {
		v, e := util.VersionString()
		if e != nil {
			panic(e)
		}
		fmt.Println(v)
	},
}
