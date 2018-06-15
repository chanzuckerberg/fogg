package init

import (
	"encoding/json"

	"github.com/ryanking/fogg/config"
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
		return e
	}
	configFile, e2 := fs.Create("fogg.json")
	if e2 != nil {
		return e2
	}
	_, e3 := configFile.Write(json)
	if e3 != nil {
		return e3
	}
	return nil
}

func Init(fs afero.Fs) error {
	project, region, bucket, profile, owner := userPrompt()
	config := config.InitConfig(project, region, bucket, profile, owner)
	e := writeConfig(fs, config)
	return e
}
