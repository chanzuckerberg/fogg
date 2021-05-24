package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	awsConfigCmd.Flags().StringP("role", "r", "default", "Use this to override the default assume role.")
	awsConfigCmd.Flags().StringP("okta-aws-saml-url", "o", "", "Your aws-okta aws_saml_url value.")
	awsConfigCmd.Flags().StringP("aws-okta-mfa-device", "m", "", "Your aws-okta mfa device to use.")

	awsConfigCmd.Flags().String("fogg-config", "fogg.yml", "Use this to override the fogg config file.")
	awsConfigCmd.Flags().String("aws-config", "~/.aws/config", "Path to your AWS config")

	expCmd.AddCommand(awsConfigCmd)
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

		var oktaAwsSamlURL *string
		if cmd.Flags().Changed("okta-aws-saml-url") {
			s, err := cmd.Flags().GetString("okta-aws-saml-url")
			if err != nil {
				return errs.WrapInternal(err, "couldn't parse okta-aws-saml-url")
			}
			oktaAwsSamlURL = &s
		}
		var awsOktaMFADevice *string
		if cmd.Flags().Changed("aws-okta-mfa-device") {
			s, err := cmd.Flags().GetString("aws-okta-mfa-device")
			if err != nil {
				return errs.WrapInternal(err, "couldn't parse aws-okta-mfa-device")
			}
			awsOktaMFADevice = &s
		}

		if (oktaAwsSamlURL != nil) != (awsOktaMFADevice != nil) {
			return errs.NewUser("Both okta-aws-saml-url and aws-okta-mfa-device must be set if either of them is.")
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
			section.Key("region").SetValue(*region)
			section.Key("output").SetValue("json")

			if oktaAwsSamlURL != nil {
				oktaProfileName := fmt.Sprintf("okta-%s", name)
				oktaSection := confIni.Section(fmt.Sprintf("profile %s", oktaProfileName))
				oktaSection.Key("role_arn").SetValue(roleARN.String())
				oktaSection.Key("aws_saml_url").SetValue(*oktaAwsSamlURL)

				// HACK HACK (el): botocore will consume both stderr and stdout
				// we need a way to display mfa prompts to users
				// direct stderr to the tty directly.
				// I assume there are some corner cases where this might break
				// but seems like a good enough place to start.
				// https://github.com/boto/botocore/issues/1348
				section.Key("credential_process").SetValue(
					fmt.Sprintf("sh -c 'aws-okta cred-process %s --mfa-duo-device %s 2> /dev/tty'", oktaProfileName, *awsOktaMFADevice))
			}
		}

		// Create directories in the aws config path if they doesn't exist already.
		dirName := filepath.Dir(awsConfigPath)
		err = os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			return err
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
