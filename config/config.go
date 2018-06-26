package config

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/spf13/afero"
)

type defaults struct {
	AWSProfileBackend  string   `json:"aws_profile_backend,omitempty"`
	AWSProfileProvider string   `json:"aws_profile_provider,omitempty"`
	AWSProviderVersion string   `json:"aws_provider_version,omitempty"`
	AWSRegion          string   `json:"aws_region"`
	AWSRegions         []string `json:"aws_regions,omitempty"`
	InfraBucket        string   `json:"infra_s3_bucket"`
	Owner              string   `json:"owner"`
	Project            string   `json:"project"`
	TerraformVersion   string   `json:"terraform_version"`
}

type Account struct {
	AccountID          *int64   `json:"account_id"`
	AWSProfileBackend  *string  `json:"aws_profile_backend"`
	AWSProfileProvider *string  `json:"aws_profile_provider"`
	AWSProviderVersion *string  `json:"aws_provider_version,omitempty"`
	AWSRegion          *string  `json:"aws_region"`  // maybe rename to provider region
	AWSRegions         []string `json:"aws_regions"` // maybe rename to provider region
	InfraBucket        *string  `json:"infra_s3_bucket"`
	Owner              *string  `json:"owner"`
	Project            *string  `json:"project"`
	TerraformVersion   *string  `json:"terraform_version"`
}

type Env struct {
	AccountID          *int64   `json:"account_id"`
	AWSProfileBackend  *string  `json:"aws_profile_backend"`
	AWSProfileProvider *string  `json:"aws_profile_provider"`
	AWSProviderVersion *string  `json:"aws_provider_version,omitempty"`
	AWSRegion          *string  `json:"aws_region"`  // maybe rename to provider region
	AWSRegions         []string `json:"aws_regions"` // maybe rename to provider region
	InfraBucket        *string  `json:"infra_s3_bucket"`
	Owner              *string  `json:"owner"`
	Project            *string  `json:"project"`
	TerraformVersion   *string  `json:"terraform_version"`

	Components map[string]*Component `json:"components"`
}

type Component struct {
	AccountID          *int64   `json:"account_id"`
	AWSProfileBackend  *string  `json:"aws_profile_backend"`
	AWSProfileProvider *string  `json:"aws_profile_provider"`
	AWSProviderVersion *string  `json:"aws_provider_version,omitempty"`
	AWSRegion          *string  `json:"aws_region"`  // maybe rename to provider region
	AWSRegions         []string `json:"aws_regions"` // maybe rename to provider region
	InfraBucket        *string  `json:"infra_s3_bucket"`
	Owner              *string  `json:"owner"`
	Project            *string  `json:"project"`
	TerraformVersion   *string  `json:"terraform_version"`
}

type Module struct {
	TerraformVersion *string `json:"terraform_version"`
}

type Config struct {
	Defaults defaults           `json:"defaults"`
	Accounts map[string]Account `json:"accounts"`
	Envs     map[string]Env     `json:"envs"`
	Modules  map[string]Module  `json:"modules"`
}

var allRegions = []string{
	"ap-south-1",
	"eu-west-3",
	"eu-west-2",
	"eu-west-1",
	"ap-northeast-2",
	"ap-northeast-1",
	"sa-east-1",
	"ca-central-1",
	"ap-southeast-1",
	"ap-southeast-2",
	"eu-central-1",
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
}

func InitConfig(project, region, bucket, awsProfile, owner string) *Config {
	return &Config{
		Defaults: defaults{
			AWSProfileBackend:  awsProfile,
			AWSProfileProvider: awsProfile,
			AWSRegion:          region,
			AWSRegions:         allRegions,
			InfraBucket:        bucket,
			Owner:              owner,
			Project:            project,
			TerraformVersion:   "0.11.0",
		},
		Accounts: map[string]Account{},
		Envs:     map[string]Env{},
		Modules:  map[string]Module{},
	}
}

func ReadConfig(f io.ReadCloser) (*Config, error) {
	defer f.Close()
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
