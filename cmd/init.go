package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	fogg_init "github.com/chanzuckerberg/fogg/init"
	prompt "github.com/segmentio/go-prompt"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type FoggProject = fogg_init.FoggProject

func init() {
	initCmd.Flags().String("project", "", "Use this to pass the project name via CLI.")
	initCmd.Flags().String("region", "", "Use this to pass the aws region via CLI.")
	initCmd.Flags().String("bucket", "", "Use this to pass the infra bucket name via CLI.")
	initCmd.Flags().String("table", "", "Use this to pass the infra dynamo table name via CLI.")
	initCmd.Flags().String("profile", "", "Use this to pass the aws auth profile via CLI.")
	initCmd.Flags().String("owner", "", "Use this to pass the owner name via CLI.")
	rootCmd.AddCommand(initCmd)
}

func isFlagPassed(cmd *cobra.Command, name string) bool {
	found := false
	cmd.Flags().Visit(func(f *pflag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func userPrompt(cmd *cobra.Command) (*FoggProject, error) {
	foggProject := new(FoggProject)
	var err error

	foggProject.Project, err = cmd.Flags().GetString("project")
	if err != nil {
		return nil, err
	}
	if foggProject.Project == "" {
		foggProject.Project = prompt.StringRequired("project name?")
	}

	foggProject.Region, err = cmd.Flags().GetString("region")
	if err != nil {
		return nil, err
	}
	if foggProject.Region == "" {
		foggProject.Region = prompt.StringRequired("aws region?")
	}

	foggProject.Bucket, err = cmd.Flags().GetString("bucket")
	if err != nil {
		return nil, err
	}
	if foggProject.Bucket == "" {
		foggProject.Bucket = prompt.StringRequired("infra bucket name?")
	}

	foggProject.Table, err = cmd.Flags().GetString("table")
	if err != nil {
		return nil, err
	}
	// check whether the flag was passed for table because table isn't required
	// so this allows passing empty string to bypass the user prompt
	if foggProject.Table == "" && !isFlagPassed(cmd, "table") {
		foggProject.Table = prompt.String("infra dynamo table name?")
	}

	foggProject.Profile, err = cmd.Flags().GetString("profile")
	if err != nil {
		return nil, err
	}
	if foggProject.Profile == "" {
		foggProject.Profile = prompt.StringRequired("auth profile?")
	}

	foggProject.Owner, err = cmd.Flags().GetString("owner")
	if err != nil {
		return nil, err
	}
	if foggProject.Owner == "" {
		foggProject.Owner = prompt.StringRequired("owner?")
	}

	return foggProject, nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new repo for use with fogg",
	Long:  "fogg init will ask you some questions and generate a basic fogg.yml.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var e error
		pwd, e := os.Getwd()
		if e != nil {
			return errs.WrapUser(e, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
		// check that we are at root of initialized git repo
		openGitOrExit(fs)

		foggProject, err := userPrompt(cmd)
		if err != nil {
			return err
		}

		return fogg_init.Init(fs, foggProject)
	},
}
