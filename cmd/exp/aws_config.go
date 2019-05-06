package exp

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/pkg/errors"
	"github.com/segmentio/go-prompt"
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
		pwd, err := os.Getwd()
		if err != nil {
			return errs.WrapUser(err, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
		// handle flags
		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			return errs.WrapInternal(err, "couldn't parse config flag")
		}
		sourceProfile, err := cmd.Flags().GetString("source-profile")
		if err != nil {
			return errs.WrapInternal(err, "couldn't parse source-profile")
		}
		role, err := cmd.Flags().GetString("role")
		if err != nil {
			return errs.WrapInternal(err, "couldn't parse role")
		}

		conf, err := config.FindAndReadConfig(fs, configFile)
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
		all := false
		choices := []string{"yes", "no", "all"}

	Loop:
		for name, account := range conf.Accounts {
			fmt.Printf("Generating config for %s\n", name)

			region := conf.Defaults.Providers.AWS.Region
			if account.Providers.AWS.Region != nil {
				region = account.Providers.AWS.Region
			}

			roleARN := arn.ARN{
				Partition: "aws",
				Service:   "iam",
				AccountID: strconv.Itoa(int(*account.Providers.AWS.AccountID)),
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

			if !all {
				fmt.Println(awsConfig.String())

				choiceIdx := prompt.Choose("Add this config?", choices)
				switch choices[choiceIdx] {
				case "no":
					continue Loop
				case "all":
					all = true
				}
			}

			err = awsConfigure(name, roleARN.String(), sourceProfile, *region)
			if err != nil {
				return err
			}
			awsConfig.Reset()
		}
		return nil
	},
}

func awsConfigure(name, roleARN, sourceProfile, region string) error {
	cmds := []struct {
		property string
		value    string
	}{
		{"role_arn", roleARN},
		{"source_profile", sourceProfile},
		{"region", region},
		{"output", "json"},
	}

	for _, params := range cmds {
		cmd := exec.Command("aws", "configure", "set", fmt.Sprintf("profile.%s.%s", name, params.property), params.value)
		err := cmd.Run()
		if err != nil {
			return errors.Wrapf(err, "Error executing: %s %s", cmd.Path, strings.Join(cmd.Args, " "))
		}
	}
	return nil
}
