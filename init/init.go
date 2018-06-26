package init

import (
	"encoding/json"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/pkg/errors"
	prompt "github.com/segmentio/go-prompt"
	"github.com/spf13/afero"
)

func userPrompt() (string, string, string, string, string) {
	project := prompt.StringRequired("project name? ")
	region := prompt.StringRequired("aws region? ")
	bucket := prompt.StringRequired("infra bucket name? ")
	profile := prompt.StringRequired("auth profile? ")
	owner := prompt.StringRequired("owner? ")

	return project, region, bucket, profile, owner
}

func writeConfig(fs afero.Fs, config *config.Config) error {
	json, e := json.MarshalIndent(config, "", "  ")
	if e != nil {
		return errors.Wrap(e, "unable to marshal json")
	}
	configFile, e := fs.Create("fogg.json")
	if e != nil {
		return errors.Wrap(e, "unable to create config file fogg.json")
	}
	_, e3 := configFile.Write(json)
	return e3
}

func Init(fs afero.Fs) error {
	project, region, bucket, profile, owner := userPrompt()
	config := config.InitConfig(project, region, bucket, profile, owner)
	e := writeConfig(fs, config)
	return e
}
