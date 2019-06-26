package cmd

import (
	"github.com/chanzuckerberg/fogg/setup"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	setupCmd.Flags().StringP("config", "c", "fogg.yml", "Use this to override the fogg config file.")
	rootCmd.AddCommand(setupCmd)
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
		logrus.Debug("setup")
		return setup.Setup(fs, config)
	},
}
