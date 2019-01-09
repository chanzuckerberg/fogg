package cmd

import (
	"github.com/chanzuckerberg/fogg/setup"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	applyCmd.Flags().BoolP("verbose", "v", false, "use this to turn on verbose output")
}

var setupCmd = &cobra.Command{
	Use:           "setup",
	Short:         "Setup dependencies for curent working directory",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fs, config, err := bootstrapCmd(cmd, debug)

		if err != nil {
			return err
		}

		// check that we are at root of initialized git repo
		openGitOrExit(fs)

		return setup.Setup(fs, config)
	},
}
