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
	initCmd.Flags().String("aws-account-id", "", "Use this to pass the primary AWS Account ID via CLI.")
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

	project, err := cmd.Flags().GetString("project")
	if err != nil {
		return nil, err
	}
	if project == "" {
		project = prompt.StringRequired("project name?")
	}
	foggProject.Project = &project

	region, err := cmd.Flags().GetString("region")
	if err != nil {
		return nil, err
	}
	if region == "" {
		region = prompt.StringRequired("aws region?")
	}
	foggProject.Region = &region

	bucket, err := cmd.Flags().GetString("bucket")
	if err != nil {
		return nil, err
	}
	// check whether the bucket flag was passed
	// bucket isn't required so this allows passing empty string to bypass the user prompt
	if bucket == "" && !isFlagPassed(cmd, "bucket") {
		bucket = prompt.String("infra bucket name?")
	}
	if bucket != "" {
		foggProject.Bucket = &bucket
	}

	table, err := cmd.Flags().GetString("table")
	if err != nil {
		return nil, err
	}
	// check whether the table flag was passed
	// table isn't required so this allows passing empty string to bypass the user prompt
	if table == "" && !isFlagPassed(cmd, "table") {
		table = prompt.String("infra dynamo table name?")
	}
	if table != "" {
		foggProject.Table = &table
	}

	profile, err := cmd.Flags().GetString("profile")
	if err != nil {
		return nil, err
	}
	// check whether the profile flag was passed
	// profile isn't required so this allows passing empty string to bypass the user prompt
	if profile == "" && !isFlagPassed(cmd, "profile") {
		profile = prompt.String("auth profile?")
	}
	if profile != "" {
		foggProject.Profile = &profile
	}

	owner, err := cmd.Flags().GetString("owner")
	if err != nil {
		return nil, err
	}
	if owner == "" {
		owner = prompt.StringRequired("owner?")
	}
	foggProject.Owner = &owner

	awsAccountID, err := cmd.Flags().GetString("aws-account-id")
	if err != nil {
		return nil, err
	}
	if awsAccountID == "" {
		awsAccountID = prompt.StringRequired("AWS Account ID?")
	}
	foggProject.AwsAccountID = &awsAccountID

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
