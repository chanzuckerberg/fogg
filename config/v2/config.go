package v2

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
)

func ReadConfig(f io.Reader) (*Config, error) {
	c := &Config{
		Docker: false,
	}
	b, e := ioutil.ReadAll(f)
	if e != nil {
		return nil, errs.WrapUser(e, "unable to read config")
	}
	e = json.Unmarshal(b, c)

	return c, errs.WrapUser(e, "unable to parse config")
}

type Config struct {
	Accounts map[string]Account   `json:"accounts,omitempty"`
	Defaults Defaults             `json:"defaults" validate:"required"`
	Docker   bool                 `json:"docker,omitempty"`
	Envs     map[string]Env       `json:"envs,omitempty"`
	Modules  map[string]v1.Module `json:"modules,omitempty"`
	Plugins  v1.Plugins           `json:"plugins,omitempty"`
	Tools    Tools                `json:"tools,omitempty"`
	Version  int                  `json:"version" validate:"required,eq=2"`
}

type common struct {
	Backend          Backend           `json:"backend,omitempty"`
	ExtraVars        map[string]string `json:"extra_vars,omitempty"`
	Owner            string            `json:"owner,omitempty" `
	Project          string            `json:"project,omitempty" `
	Providers        Providers         `json:"providers,omitempty" `
	TerraformVersion string            `json:"terraform_version,omitempty"`
}

type Defaults struct {
	common
}

type Account struct {
	common
}

type Tools struct {
	TravisCI v1.TravisCI `json:"travis_ci,omitempty"`
	TfLint   v1.TfLint   `json:"tflint,omitempty"`
}

type Env struct {
	common

	Components map[string]Component `json:"components"`
}

type Component struct {
	common

	EKS          *v1.EKSConfig     `json:"eks,omitempty"`
	Kind         *v1.ComponentKind `json:"kind,omitempty"`
	ModuleSource *string           `json:"module_source"`
}

type Providers struct {
	AWS AWSProvider `json:"aws"`
}

type AWSProvider struct {
	AccountID         *int64   `json:"account_id"`
	AdditionalRegions []string `json:"additional_regions"`
	Profile           *string  `json:"profile"`
	Region            *string  `json:"region"`
	Version           *string  `json:"version,omitempty"`
}

type Backend struct {
	Bucket      string `json:"bucket,omitempty"`
	DynamoTable string `json:"dynamodb_table,omitempty"`
	Profile     string `json:"profile,omitempty"`
	Region      string `json:"region,omitempty"`
}
