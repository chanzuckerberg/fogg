package init

import (
	"encoding/json"

	"github.com/chanzuckerberg/fogg/config"
	prompt "github.com/segmentio/go-prompt"
	"github.com/spf13/afero"
)

const DefaultPath = "fogg.json"

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
		return e
	}
	configFile, e2 := fs.Create(DefaultPath)
	if e2 != nil {
		return e2
	}
	_, e3 := configFile.Write(json)
	return e3
}

func Init(fs afero.Fs, project, region, bucket, profile, owner string) error {
	config := config.InitConfig(project, region, bucket, profile, owner)
	return writeConfig(fs, config)
}

func PromptAndInit(fs afero.Fs) error {
	project, region, bucket, profile, owner := userPrompt()
	e := Init(fs, project, region, bucket, profile, owner)
	return e
}
