package exp

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	awsConfigCmd.Flags().StringP("source-profile", "p", "default", "Use this to override the base aws profile.")
	awsConfigCmd.Flags().StringP("role", "r", "default", "Use this to override the default assume role.")
	awsConfigCmd.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")

	ExpCmd.AddCommand(awsConfigCmd)
}

var awsConfigCmd = &cobra.Command{
	Use:   "aws-config",
	Short: "Generates an ~/.aws/config from your fogg.json",
	Long:  "This command will help generate a ~/.aws/config from your fogg.json",
	RunE: func(cmd *cobra.Command, args []string) error {

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
		sourceProfile, e := cmd.Flags().GetString("source-profile")
		if e != nil {
			return errs.WrapInternal(e, "couldn't parse source-profile")
		}
		role, e := cmd.Flags().GetString("role")
		if e != nil {
			return errs.WrapInternal(e, "couldn't parse role")
		}

		config, err := config.FindAndReadConfig(fs, configFile)
		if err != nil {
			return err
		}

		templateString := `
[profile {{.accountName}}]
role_arn = {{.roleARN}}
source_profile = {{.sourceProfile}}
region = {{.region}}
output = json
`
		awsConfig := bytes.NewBuffer(nil)

		for name, account := range config.Accounts {
			region := config.Defaults.AWSRegionProvider
			if account.AWSRegionProvider != nil {
				region = *account.AWSRegionProvider
			}

			roleARN := arn.ARN{
				Partition: "aws",
				Service:   "iam",
				Region:    region,
				AccountID: strconv.Itoa(int(*account.AccountID)),
				Resource:  fmt.Sprintf("role/%s", role),
			}

			data := map[string]interface{}{
				"accountName":   name,
				"roleARN":       roleARN.String(),
				"sourceProfile": sourceProfile,
				"region":        region,
			}

			t, err := template.New("aws config").Parse(templateString)
			if err != nil {
				return errors.Wrap(err, "Could not parse template")
			}

			err = t.Execute(awsConfig, data)
			if err != nil {
				return errors.Wrap(err, "Could not templetize")
			}
		}
		// print the config
		fmt.Println(awsConfig.String())
		return nil
	},
}
