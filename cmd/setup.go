package cmd

import "github.com/spf13/cobra"

var setupCmd = &cobra.Command{
	Use:           "setup",
	Short:         "Setup dependencies for curent working directory",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		setupDebug(debug)

		return nil
	},
}
