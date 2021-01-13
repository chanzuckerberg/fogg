package cmd

import (
	"fmt"
	"os"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plan"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	componentsCmd.AddCommand(componentsListCmd)
	componentsListCmd.Flags().StringP("config", "c", "fogg.yml", "Use this to override the fogg config file.")
}

var componentsListCmd = &cobra.Command{
	Use:           "paths",
	Short:         "List paths for all configured components.",
	SilenceErrors: true, // If we don't silence here, cobra will print them. But we want to do that in cmd/root.go
	RunE: func(cmd *cobra.Command, args []string) error {
		var e error
		// Set up fs
		pwd, e := os.Getwd()
		if e != nil {
			return errs.WrapUser(e, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		// handle flags
		configFile, e := cmd.Flags().GetString("config")
		if e != nil {
			return errs.WrapInternal(e, "couldn't parse config flag")
		}

		// check that we are at root of initialized git repo
		openGitOrExit(fs)

		config, err := config.FindAndReadConfig(fs, configFile)
		if err != nil {
			return errs.WrapUser(err, "unable to parse config file")
		}

		p, e := plan.Eval(config)
		if e != nil {
			return e
		}

		fmt.Println("terraform/global")

		for _, a := range p.Accounts {
			fmt.Printf("terraform/%s\n", a.AccountName)
		}

		for _, e := range p.Envs {
			for _, c := range e.Components {
				fmt.Printf("terraform/%s/%s\n", e.Env, c.Component)
			}
		}
		return nil
	},
}
