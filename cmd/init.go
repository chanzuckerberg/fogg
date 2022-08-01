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

func userPrompt(cmd *cobra.Command) (string, string, string, string, string, string, error) {
	project, region, bucket, table, profile, owner := "", "", "", "", "", ""

	project, e := cmd.Flags().GetString("project")
	if e != nil {
		return project, region, bucket, table, profile, owner, e
	}
	if project == "" {
		project = prompt.StringRequired("project name?")
	}

	region, e = cmd.Flags().GetString("region")
	if e != nil {
		return project, region, bucket, table, profile, owner, e
	}
	if region == "" {
		region = prompt.StringRequired("aws region?")
	}

	bucket, e = cmd.Flags().GetString("bucket")
	if e != nil {
		return project, region, bucket, table, profile, owner, e
	}
	if bucket == "" {
		bucket = prompt.StringRequired("infra bucket name?")
	}

	table, e = cmd.Flags().GetString("table")
	if e != nil {
		return project, region, bucket, table, profile, owner, e
	}
	// check whether the flag was passed for table because table isn't required
	// so this allows passing empty string to bypass the user prompt
	if table == "" && !isFlagPassed(cmd, "table") {
		table = prompt.String("infra dynamo table name?")
	}

	profile, e = cmd.Flags().GetString("profile")
	if e != nil {
		return project, region, bucket, table, profile, owner, e
	}
	if profile == "" {
		profile = prompt.StringRequired("auth profile?")
	}

	owner, e = cmd.Flags().GetString("owner")
	if e != nil {
		return project, region, bucket, table, profile, owner, e
	}
	if owner == "" {
		owner = prompt.StringRequired("owner?")
	}

	return project, region, bucket, table, profile, owner, nil
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

		project, region, bucket, table, profile, owner, err := userPrompt(cmd)
		if err != nil {
			return err
		}

		return fogg_init.Init(fs, project, region, bucket, table, profile, owner)
	},
}
