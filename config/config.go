package config

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/spf13/afero"
)

type defaults struct {
	AWSProfile         string    `json:"aws_profile"`
	AWSProfileBackend  *string   `json:"aws_profile_backend,omitempty"`
	AWSProfileProvider *string   `json:"aws_profile_provider,omitempty"`
	AWSRegion          string    `json:"aws_region"`
	AWSRegions         *[]string `json:"aws_regions,omitempty"`
	InfraBucket        string    `json:"infra_s3_bucket"`
	Project            string    `json:"project"`
	TerraformVersion   string    `json:"terraform_version"`
	Owner              string    `json:"owner"`
}

type Account struct {
	AccountId          *int64    `json:"account_id"`
	AWSProfile         *string   `json:"aws_profile"`
	AWSProfileBackend  *string   `json:"aws_profile_backend"`
	AWSProfileProvider *string   `json:"aws_profile_provider"`
	AWSRegion          *string   `json:"aws_region"`  // maybe rename to provider region
	AWSRegions         *[]string `json:"aws_regions"` // maybe rename to provider region
	TerraformVersion   *string   `json:"terraform_version"`
	InfraBucket        *string   `json:"infra_s3_bucket"`
	Owner              *string   `json:"owner"`
	Project            *string   `json:"project"`
}

type Env struct{}
type Module struct {
	TerraformVersion *string `json:"terraform_version"`
}

type Config struct {
	Defaults defaults           `json:"defaults"`
	Accounts map[string]Account `json:"accounts"`
	Envs     map[string]Env     `json:"envs"`
	Modules  map[string]Module  `json:"modules"`
}

func InitConfig(project, region, bucket, aws_profile, owner string) *Config {
	return &Config{
		Defaults: defaults{
			TerraformVersion: "0.11.0",
			Project:          project,
			AWSRegion:        region,
			InfraBucket:      bucket,
			AWSProfile:       aws_profile,
			Owner:            owner,
		},
		Accounts: map[string]Account{},
		Envs:     map[string]Env{},
		Modules:  map[string]Module{},
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

func FindAndReadConfig(fs afero.Fs, configFile string) (*Config, error) {
	f, err := fs.Open(configFile)
	if err != nil {
		return nil, err
	}
	c, err2 := ReadConfig(io.ReadCloser(f))
	return c, err2
}
