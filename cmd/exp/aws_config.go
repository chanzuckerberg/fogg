package exp

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

func init() {
	awsConfigCmd.Flags().StringP("source-profile", "p", "default", "Use this to override the base aws profile.")
	awsConfigCmd.Flags().StringP("role", "r", "default", "Use this to override the default assume role.")

	awsConfigCmd.Flags().String("fogg-config", "fogg.yml", "Use this to override the fogg config file.")
	awsConfigCmd.Flags().String("aws-config", "~/.aws/config", "Path to your AWS config")
	awsConfigCmd.Flags().String("okta-aws-saml-url", "", "Your aws-okta aws_saml_url value.")
	awsConfigCmd.Flags().String("okta-role-arn", "", "Your aws-okta initial role to assume.")

	ExpCmd.AddCommand(awsConfigCmd)
}

var awsConfigCmd = &cobra.Command{
	Use:   "aws-config",
	Short: "Generates an ~/.aws/config from your fogg.yml",
	Long:  "This command will help generate a ~/.aws/config from your fogg.yml",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Set up fs
		pwd, err := os.Getwd()
		if err != nil {
			return errs.WrapUser(err, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		// handle flags
		configFile, err := cmd.Flags().GetString("fogg-config")
		if err != nil {
			return errs.WrapInternal(err, "couldn't parse fogg-config flag")
		}

		awsConfigPath, err := cmd.Flags().GetString("aws-config")
		if err != nil {
			return errs.WrapInternal(err, "couldn't parse aws-config flag")
		}
		awsConfigPath, err = homedir.Expand(awsConfigPath)
		if err != nil {
			return errs.WrapInternal(err, "couldn't expand aws-config path.")
		}

		sourceProfile, err := cmd.Flags().GetString("source-profile")
		if err != nil {
			return errs.WrapInternal(err, "couldn't parse source-profile")
		}

		var oktaAwsSamlURL *string
		if cmd.Flags().Changed("okta-aws-saml-url") {
			s, err := cmd.Flags().GetString("okta-aws-saml-url")
			if err != nil {
				return errs.WrapInternal(err, "couldn't parse okta-aws-saml-url")
			}
			oktaAwsSamlURL = &s
		}

		var oktaRoleARN *string
		if cmd.Flags().Changed("okta-role-arn") {
			s, err := cmd.Flags().GetString("okta-role-arn")
			if err != nil {
				return errs.WrapInternal(err, "couldn't parse okta-role-arn")
			}
			oktaRoleARN = &s
		}

		// both of them must be set or unset
		if (oktaAwsSamlURL != nil) != (oktaRoleARN != nil) {
			return errs.NewUser("Both okta-aws-saml-url and okta-role-arn must be set if either of them is.")
		}

		role, err := cmd.Flags().GetString("role")
		if err != nil {
			return errs.WrapInternal(err, "couldn't parse role")
		}

		conf, err := config.FindAndReadConfig(fs, configFile)
		if err != nil {
			return err
		}

		confIni := ini.Empty()

		if oktaAwsSamlURL != nil {
			section := confIni.Section("okta")
			section.Key("role_arn").SetValue(*oktaRoleARN)
			section.Key("aws_saml_url").SetValue(*oktaAwsSamlURL)
		}

		for name, account := range conf.Accounts {
			logrus.Infof("Processing profile %s", name)
			// No AWS provider, skip this account
			if account.Providers == nil || account.Providers.AWS == nil {
				logrus.Infof("Skipping %s because no AWS providers detected", name)
				continue
			}

			region := conf.Defaults.Providers.AWS.Region
			if account.Providers.AWS.Region != nil {
				region = account.Providers.AWS.Region
			}

			roleARN := arn.ARN{
				Partition: "aws",
				Service:   "iam",
				AccountID: account.Providers.AWS.AccountID.String(),
				Resource:  fmt.Sprintf("role/%s", role),
			}

			section := confIni.Section(fmt.Sprintf("profile %s", name))

			section.Key("role_arn").SetValue(roleARN.String())
			section.Key("region").SetValue(*region)
			section.Key("output").SetValue("json")

			if oktaAwsSamlURL != nil {
				section.Key("credential_source").SetValue("Environment")
			} else {
				section.Key("source_profile").SetValue(sourceProfile)
			}

		}
		awsConfigFile, err := os.OpenFile(awsConfigPath, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		defer awsConfigFile.Close()

		_, err = confIni.WriteTo(awsConfigFile)
		return err
	},
}
