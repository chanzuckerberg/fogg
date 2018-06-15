package cmd

import (
	"fmt"

	"github.com/ryanking/fogg/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of fogg",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(util.VersionString())
	},
}
