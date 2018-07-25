package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/go-playground/validator.v9"
)

type defaults struct {
	AccountID          *int64   `json:"account_id,omitempty"`
	AWSProfileBackend  string   `json:"aws_profile_backend" validate:"required"`
	AWSProfileProvider string   `json:"aws_profile_provider" validate:"required"`
	AWSProviderVersion string   `json:"aws_provider_version" validate:"required"`
	AWSRegionBackend   string   `json:"aws_region_backend" validate:"required"`
	AWSRegionProvider  string   `json:"aws_region_provider" validate:"required"`
	AWSRegions         []string `json:"aws_regions,omitempty"`
	InfraBucket        string   `json:"infra_s3_bucket" validate:"required"`
	Owner              string   `json:"owner" validate:"required"`
	Project            string   `json:"project" validate:"required"`
	SharedInfraVersion string   `json:"shared_infra_version" validate:"required"`
	TerraformVersion   string   `json:"terraform_version" validate:"required"`
}

type Account struct {
	AccountID          *int64   `json:"account_id"`
	AWSProfileBackend  *string  `json:"aws_profile_backend"`
	AWSProfileProvider *string  `json:"aws_profile_provider"`
	AWSProviderVersion *string  `json:"aws_provider_version,omitempty"`
	AWSRegionBackend   *string  `json:"aws_region_backend"`
	AWSRegionProvider  *string  `json:"aws_region_provider"`
	AWSRegions         []string `json:"aws_regions"`
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
	AWSRegionBackend   *string  `json:"aws_region_backend"`
	AWSRegionProvider  *string  `json:"aws_region_provider"`
	AWSRegions         []string `json:"aws_regions"`
	InfraBucket        *string  `json:"infra_s3_bucket"`
	Owner              *string  `json:"owner"`
	Project            *string  `json:"project"`
	TerraformVersion   *string  `json:"terraform_version"`
	Type               *string  `json:"type"`

	Components map[string]*Component `json:"components"`
}

type Component struct {
	AccountID          *int64   `json:"account_id"`
	AWSProfileBackend  *string  `json:"aws_profile_backend"`
	AWSProfileProvider *string  `json:"aws_profile_provider"`
	AWSProviderVersion *string  `json:"aws_provider_version,omitempty"`
	AWSRegionBackend   *string  `json:"aws_region_backend"`
	AWSRegionProvider  *string  `json:"aws_region_provider"`
	AWSRegions         []string `json:"aws_regions"`
	InfraBucket        *string  `json:"infra_s3_bucket"`
	Owner              *string  `json:"owner"`
	Project            *string  `json:"project"`
	SharedInfraVersion *string  `json:"shared_infra_version"`
	ModuleSource       *string  `json:"module_source"`
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

func InitConfig(project, region, bucket, awsProfile, owner, sharedInfraVersion, awsProviderVersion string) *Config {
	return &Config{
		Defaults: defaults{
			AWSProfileBackend:  awsProfile,
			AWSProfileProvider: awsProfile,
			AWSRegionBackend:   region,
			AWSRegionProvider:  region,
			AWSProviderVersion: awsProviderVersion,
			InfraBucket:        bucket,
			Owner:              owner,
			Project:            project,
			TerraformVersion:   "0.11.7",
			SharedInfraVersion: sharedInfraVersion,
		},
		Accounts: map[string]Account{},
		Envs:     map[string]Env{},
		Modules:  map[string]Module{},
	}
}

func ReadConfig(f io.Reader) (*Config, error) {
	c := &Config{}
	b, e := ioutil.ReadAll(f)
	if e != nil {
		return nil, errors.Wrap(e, "unable to read config")
	}
	e = json.Unmarshal(b, c)
	if e != nil {
		return nil, errors.Wrap(e, "unable to parse json config file")
	}
	return c, nil
}

func FindAndReadConfig(fs afero.Fs, configFile string) (*Config, error) {
	f, e := fs.Open(configFile)
	if e != nil {
		return nil, errors.Wrap(e, "unable to open config file")
	}
	reader := io.ReadCloser(f)
	defer reader.Close()
	c, err2 := ReadConfig(reader)
	return c, err2
}

func (c *Config) Validate() error {
	v := validator.New()
	// https://github.com/go-playground/validator/issues/323#issuecomment-343670840
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})
	return v.Struct(c)
}
