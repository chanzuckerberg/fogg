package cmd

import (
	"fmt"

	"github.com/chanzuckerberg/fogg/apply"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/spf13/cobra"
)

var environment string
var component string

func init() {
	applyCmd.Flags().StringP("config", "c", "fogg.yml", "Use this to override the fogg config file.")
	applyCmd.Flags().BoolP("upgrade", "u", false, "Use this when running a new version of fogg")
	applyCmd.Flags().StringVarP(&environment, "env", "e", "", "Limit apply to specific environment")
	applyCmd.Flags().StringVarP(&component, "component", "f", "", "Limit apply to specific component (requires env flag)")
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

		var envFilter *string
		if environment != "" {
			envFilter = &environment
		}

		var compFilter *string
		if component != "" {
			compFilter = &component
			if envFilter == nil {
				return fmt.Errorf("component flag requires env flag")
			}
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
		e = apply.Apply(fs, config, templates.Templates, upgrade, envFilter, compFilter)

		return e
	},
}
