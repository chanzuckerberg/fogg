package cmd

import (
	"github.com/chanzuckerberg/fogg/apply"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	applyCmd.Flags().BoolP("verbose", "v", false, "use this to turn on verbose output")
	applyCmd.Flags().BoolP("upgrade", "u", false, "use this when running a new version of fogg")
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:           "apply",
	Short:         "Apply model defined in fogg.json to the current tree.",
	Long:          "This command will take the model defined in fogg.json, build a plan and generate the appropriate files from templates.",
	SilenceErrors: true, // If we don't silence here, cobra will print them. But we want to do that in cmd/root.go
	RunE: func(cmd *cobra.Command, args []string) error {
		setupDebug(debug)

		var e error
		fs, e := openFs()
		if e != nil {
			return e
		}

		// handle flags
		verbose, e := cmd.Flags().GetBool("verbose")
		if e != nil {
			return e
		}
		configFile, e := cmd.Flags().GetString("config")
		if e != nil {
			return e
		}

		upgrade, e := cmd.Flags().GetBool("upgrade")
		if e != nil {
			return e
		}

		// check that we are at root of initialized git repo
		openGitOrExit(fs)

		config, err := readAndValidateConfig(fs, configFile, verbose)

		e = mergeConfigValidationErrors(err)
		if e != nil {
			return e
		}

		// apply
		e = apply.Apply(fs, config, templates.Templates, upgrade)

		return e
	},
}
