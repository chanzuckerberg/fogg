package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/hashicorp/go-multierror"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/go-playground/validator.v9"
)

type defaults struct {
	AccountID          *int64            `json:"account_id,omitempty"`
	AWSProfileBackend  string            `json:"aws_profile_backend" validate:"required"`
	AWSProfileProvider string            `json:"aws_profile_provider" validate:"required"`
	AWSProviderVersion string            `json:"aws_provider_version" validate:"required"`
	AWSRegionBackend   string            `json:"aws_region_backend" validate:"required"`
	AWSRegionProvider  string            `json:"aws_region_provider" validate:"required"`
	AWSRegions         []string          `json:"aws_regions,omitempty"`
	ExtraVars          map[string]string `json:"extra_vars"`
	InfraBucket        string            `json:"infra_s3_bucket" validate:"required"`
	Owner              string            `json:"owner" validate:"required"`
	Project            string            `json:"project" validate:"required"`
	TerraformVersion   string            `json:"terraform_version" validate:"required"`
}

type Account struct {
	AccountID          *int64            `json:"account_id"`
	AWSProfileBackend  *string           `json:"aws_profile_backend"`
	AWSProfileProvider *string           `json:"aws_profile_provider"`
	AWSProviderVersion *string           `json:"aws_provider_version,omitempty"`
	AWSRegionBackend   *string           `json:"aws_region_backend"`
	AWSRegionProvider  *string           `json:"aws_region_provider"`
	AWSRegions         []string          `json:"aws_regions"`
	ExtraVars          map[string]string `json:"extra_vars,omitempty"`
	InfraBucket        *string           `json:"infra_s3_bucket"`
	Owner              *string           `json:"owner"`
	Project            *string           `json:"project"`
	TerraformVersion   *string           `json:"terraform_version"`
}

type Env struct {
	AccountID          *int64            `json:"account_id"`
	AWSProfileBackend  *string           `json:"aws_profile_backend"`
	AWSProfileProvider *string           `json:"aws_profile_provider"`
	AWSProviderVersion *string           `json:"aws_provider_version,omitempty"`
	AWSRegionBackend   *string           `json:"aws_region_backend"`
	AWSRegionProvider  *string           `json:"aws_region_provider"`
	AWSRegions         []string          `json:"aws_regions"`
	ExtraVars          map[string]string `json:"extra_vars,omitempty"`
	InfraBucket        *string           `json:"infra_s3_bucket"`
	Owner              *string           `json:"owner"`
	Project            *string           `json:"project"`
	TerraformVersion   *string           `json:"terraform_version"`
	Type               *string           `json:"type"`

	Components map[string]*Component `json:"components"`
}

type Component struct {
	AccountID          *int64            `json:"account_id"`
	AWSProfileBackend  *string           `json:"aws_profile_backend"`
	AWSProfileProvider *string           `json:"aws_profile_provider"`
	AWSProviderVersion *string           `json:"aws_provider_version,omitempty"`
	AWSRegionBackend   *string           `json:"aws_region_backend"`
	AWSRegionProvider  *string           `json:"aws_region_provider"`
	AWSRegions         []string          `json:"aws_regions"`
	ExtraVars          map[string]string `json:"extra_vars,omitempty"`
	InfraBucket        *string           `json:"infra_s3_bucket"`
	ModuleSource       *string           `json:"module_source"`
	Owner              *string           `json:"owner"`
	Project            *string           `json:"project"`
	TerraformVersion   *string           `json:"terraform_version"`
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

func InitConfig(project, region, bucket, awsProfile, owner, awsProviderVersion string) *Config {
	return &Config{
		Defaults: defaults{
			AWSProfileBackend:  awsProfile,
			AWSProfileProvider: awsProfile,
			AWSProviderVersion: awsProviderVersion,
			AWSRegionBackend:   region,
			AWSRegionProvider:  region,
			ExtraVars:          map[string]string{},
			InfraBucket:        bucket,
			Owner:              owner,
			Project:            project,
			TerraformVersion:   "0.11.7",
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

func (c *Config) validateExtraVars() error {
	var err *multierror.Error
	validate := func(extraVars map[string]string) {
		for extraVar := range extraVars {
			if _, ok := reservedVariableNames[extraVar]; ok {
				err = multierror.Append(err, fmt.Errorf("extra_var[%s] is a fogg reserved variable name", extraVar))
			}
		}
		return
	}

	extraVars := []map[string]string{}
	extraVars = append(extraVars, c.Defaults.ExtraVars)
	for _, env := range c.Envs {
		extraVars = append(extraVars, env.ExtraVars)
		for _, component := range env.Components {
			extraVars = append(extraVars, component.ExtraVars)
		}
	}
	for _, extraVar := range extraVars {
		validate(extraVar)
	}

	return errors.Wrap(err.ErrorOrNil(), "extra_vars contains reserved names")
}
