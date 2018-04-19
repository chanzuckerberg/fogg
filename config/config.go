package config

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/spf13/afero"
)

type defaults struct {
	AWSRegion        string `json:"aws_region"`
	AWSProfile       string `json:"aws_profile"`
	InfraBucket      string `json:"infra_bucket"`
	Project          string `json:"project"`
	SharedInfraPath  string `json:"shared_infra_base"`
	TerraformVersion string `json:"terraform_version"`
	// regions
	// shared infra version
	// owner
	// aws_profile_backend
	// aws_profile_provider
}

type account struct {
	defaults
}

type Config struct {
	Defaults defaults `json:"defaults"`
	// Envs     map[string]env     `json:"envs"`
	// Modules  map[string]module  `json:"modules"`
	Accounts map[string]account `json:"account"`
}

func DefaultConfig() *Config {
	return &Config{
		Defaults: defaults{
			SharedInfraPath:  "git@github.com:chanzuckerberg/shared-infra//",
			TerraformVersion: "0.11.0",
		},
	}
}

func ReadConfig(f io.ReadCloser) (*Config, error) {
	c := &Config{}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err2 := json.Unmarshal(b, c)
	if err2 != nil {
		return nil, err2
	}
	return c, nil
}

func FindAndReadConfig(fs afero.Fs) (*Config, error) {
	f, err := fs.Open("fogg.json")
	if err != nil {
		return nil, err
	}
	c, err2 := ReadConfig(io.ReadCloser(f))
	return c, err2
}
