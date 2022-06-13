package cmd

import (
	"github.com/chanzuckerberg/fogg/apply"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd.Flags().StringP("config", "c", "fogg.yml", "Use this to override the fogg config file.")
	applyCmd.Flags().BoolP("upgrade", "u", false, "Use this when running a new version of fogg")
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:           "apply",
	Short:         "Apply model defined in fogg.yml to the current tree.",
	Long:          "This command will take the model defined in fogg.yml, build a plan and generate the appropriate files from templates.",
	SilenceErrors: true, // If we don't silence here, cobra will print them. But we want to do that in cmd/root.go
	RunE: func(cmd *cobra.Command, args []string) error {
		var e error
		fs, e := openFs()
		if e != nil {
			return e
		}

		// handle flags

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

		config, warnings, err := readAndValidateConfig(fs, configFile)
		printWarnings(warnings)

		e = mergeConfigValidationErrors(err)
		if e != nil {
			return e
		}

		// apply
		e = apply.Apply(fs, config, templates.Templates, upgrade)

		return e
	},
}
