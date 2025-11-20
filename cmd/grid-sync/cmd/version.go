package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of grid-sync",
	RunE: func(cmd *cobra.Command, args []string) error {
		v, e := util.VersionString()
		if e != nil {
			return e
		}
		fmt.Println(v)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
