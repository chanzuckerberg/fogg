package init

import (
	"github.com/chanzuckerberg/fogg/config"
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	prompt "github.com/segmentio/go-prompt"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

const AWSProviderVersion = "2.47.0"

func userPrompt() (string, string, string, string, string, string) {
	project := prompt.StringRequired("project name?")
	region := prompt.StringRequired("aws region?")
	bucket := prompt.StringRequired("infra bucket name?")
	table := prompt.String("infra dynamo table name?")
	profile := prompt.StringRequired("auth profile?")
	owner := prompt.StringRequired("owner?")

	return project, region, bucket, table, profile, owner
}

func writeConfig(fs afero.Fs, config *v2.Config) error {
	yaml, err := yaml.Marshal(config)
	if err != nil {
		return errs.WrapInternal(err, "unable to marshal yaml")
	}

	yamlConfigFile, err := fs.Create("fogg.yml")
	if err != nil {
		return errs.WrapInternal(err, "unable to create config file fogg.yml")
	}
	_, err = yamlConfigFile.Write(yaml)
	return err
}

//Init reads user console input and generates a fogg.yml file
func Init(fs afero.Fs) error {
	project, region, bucket, table, profile, owner := userPrompt()
	config := config.InitConfig(project, region, bucket, table, profile, owner, AWSProviderVersion)
	e := writeConfig(fs, config)
	return e
}
