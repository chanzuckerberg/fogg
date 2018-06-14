package config

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/spf13/afero"
)

type defaults struct {
	AWSProfile         string   `json:"aws_profile"`
	AWSProfileBackend  *string  `json:"aws_profile_backend"`
	AWSProfileProvider *string  `json:"aws_profile_provider"`
	AWSRegion          string   `json:"aws_region"`
	AWSRegions         []string `json:"aws_regions"`
	InfraBucket        string   `json:"infra_bucket"`
	Project            string   `json:"project"`
	SharedInfraPath    string   `json:"shared_infra_base"`
	TerraformVersion   string   `json:"terraform_version"`
	// shared infra version
	// owner
	// aws_profile_backend
	// aws_profile_provider
}

type Account struct {
	AccountId          *int64    `json:"account_id"`
	AWSProfile         *string   `json:"aws_profile"`
	AWSProfileBackend  *string   `json:"aws_profile_backend"`
	AWSProfileProvider *string   `json:"aws_profile_provider"`
	AWSRegion          *string   `json:"aws_region"`  // maybe rename to provider region
	AWSRegions         *[]string `json:"aws_regions"` // maybe rename to provider region
}

type Config struct {
	Defaults defaults           `json:"defaults"`
	Accounts map[string]Account `json:"accounts"`
	// Envs     map[string]env     `json:"envs"`
	// Modules  map[string]module  `json:"modules"`
}

func DefaultConfig() *Config {
	return &Config{
		Defaults: defaults{
			SharedInfraPath:  "git@github.com:chanzuckerberg/shared-infra//",
			TerraformVersion: "0.11.0",
			AWSRegions: []string{
				"ap-northeast-1",
				"ap-northeast-2",
				"ap-south-1",
				"ap-southeast-1",
				"ap-southeast-2",
				"ca-central-1",
				"eu-central-1",
				"eu-west-1",
				"eu-west-2",
				"eu-west-3",
				"sa-east-1",
				"us-east-1",
				"us-east-2",
				"us-west-1",
				"us-west-2",
			},
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
