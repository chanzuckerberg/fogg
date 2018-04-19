package main

import (
	"encoding/json"

	"github.com/ryanking/fogg/config"
	prompt "github.com/segmentio/go-prompt"
	"github.com/spf13/afero"
)

func userPrompt() (string, string, string, string) {
	project := prompt.StringRequired("project name? ")
	region := prompt.StringRequired("aws region? ")
	bucket := prompt.StringRequired("infra bucket name? ")
	profile := prompt.StringRequired("auth profile? ")

	return project, region, bucket, profile
}

func createConfig(project, region, bucket, profile string) *config.Config {
	c := config.DefaultConfig()
	c.Defaults.Project = project
	c.Defaults.AWSRegion = region
	c.Defaults.InfraBucket = bucket
	c.Defaults.AWSProfile = profile

	return c
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
	project, region, bucket, profile := userPrompt()
	config := createConfig(project, region, bucket, profile)
	e := writeConfig(fs, config)
	return e
}
