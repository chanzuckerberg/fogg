package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plugins"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

type TfLint struct {
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

type Defaults struct {
	AccountID          int64             `json:"account_id,omitempty" yaml:"account_id,omitempty" validate:"required"`
	AWSProfileBackend  string            `json:"aws_profile_backend" yaml:"aws_profile_backend" validate:"required"`
	AWSProfileProvider string            `json:"aws_profile_provider" yaml:"aws_profile_provider" validate:"required"`
	AWSProviderVersion string            `json:"aws_provider_version" yaml:"aws_provider_version" validate:"required"`
	AWSRegionBackend   string            `json:"aws_region_backend" yaml:"aws_region_backend" validate:"required"`
	AWSRegionProvider  string            `json:"aws_region_provider" yaml:"aws_region_provider" validate:"required"`
	AWSRegions         []string          `json:"aws_regions,omitempty" yaml:"aws_regions,omitempty" `
	ExtraVars          map[string]string `json:"extra_vars" yaml:"extra_vars" `
	InfraBucket        string            `json:"infra_s3_bucket" yaml:"infra_s3_bucket" validate:"required"`
	InfraDynamoTable   string            `json:"infra_dynamo_db_table" yaml:"infra_dynamo_db_table" `
	Owner              string            `json:"owner" yaml:"owner" validate:"required"`
	Project            string            `json:"project" yaml:"project" validate:"required"`
	TerraformVersion   string            `json:"terraform_version" yaml:"terraform_version" validate:"required"`
	TfLint             *TfLint           `json:"tflint,omitempty" yaml:"tflint,omitempty"`
}

type Account struct {
	AccountID          *int64            `json:"account_id" yaml:"account_id"`
	AWSProfileBackend  *string           `json:"aws_profile_backend" yaml:"aws_profile_backend"`
	AWSProfileProvider *string           `json:"aws_profile_provider" yaml:"aws_profile_provider"`
	AWSProviderVersion *string           `json:"aws_provider_version,omitempty" yaml:"aws_provider_version,omitempty"`
	AWSRegionBackend   *string           `json:"aws_region_backend" yaml:"aws_region_backend"`
	AWSRegionProvider  *string           `json:"aws_region_provider" yaml:"aws_region_provider"`
	AWSRegions         []string          `json:"aws_regions" yaml:"aws_regions"`
	ExtraVars          map[string]string `json:"extra_vars,omitempty" yaml:"extra_vars,omitempty"`
	InfraBucket        *string           `json:"infra_s3_bucket" yaml:"infra_s3_bucket"`
	InfraDynamoTable   *string           `json:"infra_dynamo_db_table" yaml:"infra_dynamo_db_table"`
	Owner              *string           `json:"owner" yaml:"owner"`
	Project            *string           `json:"project" yaml:"project"`
	TerraformVersion   *string           `json:"terraform_version" yaml:"terraform_version"`
	TfLint             *TfLint           `json:"tflint,omitempty" yaml:"tflint,omitempty"`
}

type Env struct {
	AccountID          *int64            `json:"account_id" yaml:"account_id"`
	AWSProfileBackend  *string           `json:"aws_profile_backend" yaml:"aws_profile_backend"`
	AWSProfileProvider *string           `json:"aws_profile_provider" yaml:"aws_profile_provider"`
	AWSProviderVersion *string           `json:"aws_provider_version,omitempty" yaml:"aws_provider_version,omitempty"`
	AWSRegionBackend   *string           `json:"aws_region_backend" yaml:"aws_region_backend"`
	AWSRegionProvider  *string           `json:"aws_region_provider" yaml:"aws_region_provider"`
	AWSRegions         []string          `json:"aws_regions" yaml:"aws_regions"`
	ExtraVars          map[string]string `json:"extra_vars,omitempty" yaml:"extra_vars,omitempty"`
	InfraBucket        *string           `json:"infra_s3_bucket" yaml:"infra_s3_bucket"`
	InfraDynamoTable   *string           `json:"infra_dynamo_db_table" yaml:"infra_dynamo_db_table"`
	Owner              *string           `json:"owner" yaml:"owner"`
	Project            *string           `json:"project" yaml:"project"`
	TerraformVersion   *string           `json:"terraform_version" yaml:"terraform_version"`
	TfLint             *TfLint           `json:"tflint,omitempty" yaml:"tflint,omitempty"`

	Components map[string]*Component `json:"components" yaml:"components"`
}

// ComponentKind is the kind of this component
type ComponentKind string

// GetOrDefault gets the component kind or default
func (ck *ComponentKind) GetOrDefault() ComponentKind {
	if ck == nil || *ck == "" {
		return DefaultComponentKind
	}
	return *ck
}

const (
	// DefaultComponentKind defaults to terraform component
	DefaultComponentKind ComponentKind = "terraform"
	// ComponentKindTerraform is a terraform component
	ComponentKindTerraform = DefaultComponentKind
	// ComponentKindHelmTemplate is a helm template component
	ComponentKindHelmTemplate ComponentKind = "helm_template"
)

// EKSConfig is the configuration for an eks cluster
type EKSConfig struct {
	ClusterName string `json:"cluster_name" yaml:"cluster_name"`
}

type Component struct {
	AccountID          *int64            `json:"account_id" yaml:"account_id"`
	AWSProfileBackend  *string           `json:"aws_profile_backend" yaml:"aws_profile_backend"`
	AWSProfileProvider *string           `json:"aws_profile_provider" yaml:"aws_profile_provider"`
	AWSProviderVersion *string           `json:"aws_provider_version,omitempty" yaml:"aws_provider_version,omitempty"`
	AWSRegionBackend   *string           `json:"aws_region_backend" yaml:"aws_region_backend"`
	AWSRegionProvider  *string           `json:"aws_region_provider" yaml:"aws_region_provider"`
	AWSRegions         []string          `json:"aws_regions" yaml:"aws_regions"`
	EKS                *EKSConfig        `json:"eks,omitempty" yaml:"eks,omitempty"`
	ExtraVars          map[string]string `json:"extra_vars,omitempty" yaml:"extra_vars,omitempty"`
	InfraBucket        *string           `json:"infra_s3_bucket" yaml:"infra_s3_bucket"`
	InfraDynamoTable   *string           `json:"infra_dynamo_db_table" yaml:"infra_dynamo_db_table"`
	Kind               *ComponentKind    `json:"kind,omitempty" yaml:"kind,omitempty"`
	ModuleSource       *string           `json:"module_source" yaml:"module_source"`
	Owner              *string           `json:"owner" yaml:"owner"`
	Project            *string           `json:"project" yaml:"project"`
	TerraformVersion   *string           `json:"terraform_version" yaml:"terraform_version"`
	TfLint             *TfLint           `json:"tflint,omitempty" yaml:"tflint,omitempty"`
}

// Plugins contains configuration around plugins
type Plugins struct {
	CustomPlugins      map[string]*plugins.CustomPlugin `json:"custom_plugins,omitempty" yaml:"custom_plugins,omitempty"`
	TerraformProviders map[string]*plugins.CustomPlugin `json:"terraform_providers,omitempty" yaml:"terraform_providers,omitempty"`
}

// Module is a module
type Module struct {
	TerraformVersion *string `json:"terraform_version,omitempty" yaml:"terraform_version,omitempty"`
}

type TravisCI struct {
	Enabled        bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	AWSIAMRoleName string `json:"aws_iam_role_name" yaml:"aws_iam_role_name"`
	TestBuckets    int    `json:"test_buckets" yaml:"test_buckets"`
}

type Config struct {
	Accounts map[string]Account `json:"accounts" yaml:"accounts"`
	Defaults Defaults           `json:"defaults" yaml:"defaults"`
	Docker   bool               `json:"docker,omitempty" yaml:"docker,omitempty"`
	Envs     map[string]Env     `json:"envs" yaml:"envs"`
	Modules  map[string]Module  `json:"modules" yaml:"modules"`
	Plugins  Plugins            `json:"plugins,omitempty" yaml:"plugins"`
	TravisCI *TravisCI          `json:"travis_ci,omitempty" yaml:"travis_ci,omitempty"`
}

func ReadConfig(b []byte) (*Config, error) {
	var e error
	c := &Config{
		Docker: true,
	}

	if IsJSON(b) {
		e = json.Unmarshal(b, c)
		logrus.Warn("JSON is deprecated, consider migrating to yaml")
	} else {
		e = yaml.Unmarshal(b, c)
	}
	if e != nil {
		return nil, errs.WrapUser(e, "unable to parse config file")
	}

	return c, errs.WrapUser(e, "unable to parse yaml config file")
}

func IsJSON(b []byte) bool {
	jsonPrefix := []byte("{")
	trimmed := bytes.TrimLeftFunc(b, unicode.IsSpace)

	return bytes.HasPrefix(trimmed, jsonPrefix)
}

// Validate validates the config
func (c *Config) Validate() error {
	err := c.validateExtraVars()
	if err != nil {
		return err
	}

	v := validator.New()
	// https://github.com/go-playground/validator/issues/323#issuecomment-343670840
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("yml"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})
	return v.Struct(c)
}

// validateExtraVars make sure users don't specify reserved variables
func (c *Config) validateExtraVars() error {
	var err *multierror.Error
	validate := func(extraVars map[string]string) {
		for extraVar := range extraVars {
			if _, ok := ReservedVariableNames[extraVar]; ok {
				err = multierror.Append(err, fmt.Errorf("extra_var[%s] is a fogg reserved variable name", extraVar))
			}
		}
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

	if err.ErrorOrNil() != nil {
		return errs.WrapUser(err.ErrorOrNil(), "extra_vars contains reserved variable names")
	}
	return nil
}
