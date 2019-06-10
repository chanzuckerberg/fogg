package init

import (
	"github.com/chanzuckerberg/fogg/config"
	v1 "github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
	prompt "github.com/segmentio/go-prompt"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

const AWSProviderVersion = "1.27.0"

func userPrompt() (string, string, string, string, string, string) {
	project := prompt.StringRequired("project name?")
	region := prompt.StringRequired("aws region?")
	bucket := prompt.StringRequired("infra bucket name?")
	table := prompt.String("infra dynamo table name?")
	profile := prompt.StringRequired("auth profile?")
	owner := prompt.StringRequired("owner?")

	return project, region, bucket, table, profile, owner
}

func writeConfig(fs afero.Fs, config *v1.Config) error {
	yaml, yamlErr := yaml.Marshal(config)

	if yamlErr != nil {
		return errs.WrapInternal(yamlErr, "unable to marshal yaml")
	}

	yamlConfigFile, yamlErr := fs.Create("fogg.yml")
	if yamlErr != nil {
		return errs.WrapInternal(yamlErr, "unable to create config file fogg.yml")
	}
	_, yamlStatus := yamlConfigFile.Write(yaml)

	return yamlStatus
}

//Init reads user console input and generates a fogg.yaml file
func Init(fs afero.Fs) error {
	project, region, bucket, table, profile, owner := userPrompt()
	config := config.InitConfig(project, region, bucket, table, profile, owner, AWSProviderVersion)
	e := writeConfig(fs, config)
	return e
}
