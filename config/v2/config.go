package v2

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"path/filepath"
	"reflect"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plugins"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

//ReadConfig take a byte array as input and outputs a json or yaml config struct
func ReadConfig(fs afero.Fs, b []byte, configFile string) (*Config, error) {
	var e error
	c := &Config{
		Docker: false,
	}

	info, e := fs.Stat(configFile)
	if e != nil {
		return nil, errs.WrapUser(e, "unable to find file")
	}

	ext := filepath.Ext(info.Name())
	//Determines the file extension
	switch ext {
	case ".yml", ".yaml":
		reader := bytes.NewReader(b)
		decoder := yaml.NewDecoder(reader)
		decoder.KnownFields(true)
		e = decoder.Decode(c)
	default:
		return nil, errs.NewUserf("File type %s is not supported", ext)
	}
	return c, e
}

type Config struct {
	Accounts map[string]Account `yaml:"accounts,omitempty"`
	Defaults Defaults           `yaml:"defaults" validate:"required"`
	Docker   bool               `yaml:"docker,omitempty"`
	Envs     map[string]Env     `yaml:"envs,omitempty"`
	Global   Component          `yaml:"global,omitempty"`
	Modules  map[string]Module  `yaml:"modules,omitempty"`
	Plugins  Plugins            `yaml:"plugins,omitempty"`
	Version  int                `validate:"required,eq=2"`
}

type Common struct {
	Backend          *Backend          `yaml:"backend,omitempty"`
	ExtraVars        map[string]string `yaml:"extra_vars,omitempty"`
	Owner            *string           `yaml:"owner,omitempty"`
	Project          *string           `yaml:"project,omitempty"`
	Providers        *Providers        `yaml:"providers,omitempty"`
	TerraformVersion *string           `yaml:"terraform_version,omitempty"`
	Tools            *Tools            `yaml:"tools,omitempty"`
}

type Defaults struct {
	Common `yaml:",inline"`
}

type Account struct {
	Common `yaml:",inline"`
}

type Tools struct {
	TravisCI        *TravisCI        `yaml:"travis_ci,omitempty"`
	CircleCI        *CircleCI        `yaml:"circle_ci,omitempty"`
	GitHubActionsCI *GitHubActionsCI `yaml:"github_actions_ci,omitempty"`
	TfLint          *TfLint          `yaml:"tflint,omitempty"`
}

type CircleCI struct {
	CommonCI `yaml:",inline"`

	SSHKeyFingerprints []string `yaml:"ssh_key_fingerprints"`
}

type GitHubActionsCI struct {
	CommonCI `yaml:",inline"`
}

type Env struct {
	Common `yaml:",inline"`

	Components map[string]Component `yaml:"components,omitempty"`
}

type Component struct {
	Common `yaml:",inline"`

	EKS          *EKSConfig     `yaml:"eks,omitempty"`
	Kind         *ComponentKind `yaml:"kind,omitempty"`
	ModuleSource *string        `yaml:"module_source,omitempty"`
}

type Providers struct {
	AWS       *AWSProvider       `yaml:"aws,omitempty"`
	Bless     *BlessProvider     `yaml:"bless,omitempty"`
	Github    *GithubProvider    `yaml:"github,omitempty"`
	Heroku    *HerokuProvider    `yaml:"heroku,omitempty"`
	Okta      *OktaProvider      `yaml:"okta,omitempty"`
	Snowflake *SnowflakeProvider `yaml:"snowflake,omitempty"`
	Datadog   *DatadogProvider   `yaml:"datadog,omitempty"`
	Tfe       *TfeProvider       `yaml:"tfe,omitempty"`
}

// CommonProvider encapsulates common properties across providers
// TODO refactor other providers to use CommonProvider inline
type CommonProvider struct {
	Enabled *bool   `yaml:"enabled,omitempty"`
	Version *string `yaml:"version,omitempty"`
}

// OktaProvider is an okta provider
type OktaProvider struct {
	// the okta provider is optional (above) but if supplied you must set an OrgName
	OrgName *string `yaml:"org_name,omitempty"`
	Version *string `yaml:"version,omitempty"`
}

// BlessProvider allows for terraform-provider-bless configuration
type BlessProvider struct {
	// the bless provider is optional (above) but if supplied you must set a region and aws_profile
	AdditionalRegions []string `yaml:"additional_regions,omitempty"`
	AWSProfile        *string  `yaml:"aws_profile,omitempty"`
	AWSRegion         *string  `yaml:"aws_region,omitempty"`
	Version           *string  `yaml:"version,omitempty"`
}

type AWSProvider struct {
	// the aws provider is optional (above) but if supplied you must set account id and region
	AccountID         *json.Number `yaml:"account_id,omitempty"`
	AdditionalRegions []string     `yaml:"additional_regions,omitempty"`
	Profile           *string      `yaml:"profile,omitempty"`
	Region            *string      `yaml:"region,omitempty"`
	Role              *string      `yaml:"role,omitempty"` // FIXME validate format
	Version           *string      `yaml:"version,omitempty"`
}

type GithubProvider struct {
	Organization *string `yaml:"organization,omitempty"`
	BaseURL      *string `yaml:"base_url,omitempty"`
	Version      *string `yaml:"version,omitempty"`
}

type SnowflakeProvider struct {
	Account *string `yaml:"account,omitempty"`
	Role    *string `yaml:"role,omitempty"`
	Region  *string `yaml:"region,omitempty"`
	Version *string `yaml:"version,omitempty"`
}

type HerokuProvider struct {
	Version *string `yaml:"version,omitempty"`
}

type DatadogProvider struct {
	Version *string `yaml:"version,omitempty"`
}

type TfeProvider struct {
	CommonProvider `yaml:",inline"`

	Hostname *string `yaml:"hostname,omitempty"`
}

//Backend is used to configure a terraform backend
type Backend struct {
	Kind *string `yaml:"kind,omitempty" validate:"omitempty,oneof=s3 remote"`

	// fields used for S3 backend
	AccountID   *string `yaml:"account_id,omitempty"`
	Bucket      *string `yaml:"bucket,omitempty"`
	DynamoTable *string `yaml:"dynamodb_table,omitempty"`
	Profile     *string `yaml:"profile,omitempty"`
	Region      *string `yaml:"region,omitempty"`
	Role        *string `yaml:"role,omitempty"`

	// fields used for remote backend
	HostName     *string `yaml:"host_name,omitempty"`
	Organization *string `yaml:"organization,omitempty"`
}

// Module is a module
type Module struct {
	TerraformVersion *string `yaml:"terraform_version,omitempty"`
}

// Plugins contains configuration around plugins
type Plugins struct {
	CustomPlugins      map[string]*plugins.CustomPlugin `yaml:"custom_plugins,omitempty"`
	TerraformProviders map[string]*plugins.CustomPlugin `yaml:"terraform_providers,omitempty"`
}

type TravisCI struct {
	CommonCI `yaml:",inline"`
}

type CommonCI struct {
	Enabled        *bool                       `yaml:"enabled,omitempty"`
	AWSIAMRoleName *string                     `yaml:"aws_iam_role_name,omitempty"`
	TestBuckets    *int                        `yaml:"test_buckets,omitempty"`
	Command        *string                     `yaml:"command,omitempty"`
	Buildevents    *bool                       `yaml:"buildevents,omitempty"`
	Providers      map[string]CIProviderConfig `yaml:"providers,omitempty"`
}

type CIProviderConfig struct {
	Disabled bool `yaml:"disabled,omitempty"`
}

type TfLint struct {
	Enabled *bool `yaml:"enabled,omitempty"`
}

// EKSConfig is the configuration for an eks cluster
type EKSConfig struct {
	ClusterName string `yaml:"cluster_name"`
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

// Generate is used for test/quick integration. There are supposedly ways to do this without polluting the public
// api, but givent that fogg isn't an api, it doesn't seem like a big deal
func (c *Config) Generate(r *rand.Rand, size int) reflect.Value {
	// TODO write this to be part of tests https://github.com/shiwano/submarine/blob/5c02c0cfcf05126454568ef9624550eb0d84f86c/server/battle/src/battle/util/util_test.go#L19

	conf := &Config{}

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	randString := func(r *rand.Rand, n int) string {
		b := make([]byte, n)
		for i := range b {
			b[i] = letterBytes[r.Intn(len(letterBytes))]
		}
		return string(b)
	}

	randNonEmptyString := func(r *rand.Rand, s int) string {
		return "asdf"
	}

	randStringPtr := func(r *rand.Rand, s int) *string {
		str := randString(r, s)
		return &str
	}

	randBoolPtr := func(r *rand.Rand, s int) *bool {
		b := r.Float32() > 0.5
		return &b
	}
	randIntPtr := func(r *rand.Rand, s int) *int {
		i := r.Intn(s)
		return &i
	}

	randStringMap := func(r *rand.Rand, s int) map[string]string {
		m := map[string]string{}

		for i := 0; i < s; i++ {
			m[randNonEmptyString(r, s)] = randString(r, s)
		}

		return map[string]string{}
	}

	randOktaProvider := func(r *rand.Rand, s int) *OktaProvider {
		if r.Float32() < 0.5 {
			return nil
		}
		return &OktaProvider{
			Version: randStringPtr(r, s),
			OrgName: randStringPtr(r, s),
		}
	}

	randBlessProvider := func(r *rand.Rand, s int) *BlessProvider {
		if r.Float32() < 0.5 {
			return nil
		}
		return &BlessProvider{
			Version:           randStringPtr(r, s),
			AWSRegion:         randStringPtr(r, s),
			AWSProfile:        randStringPtr(r, s),
			AdditionalRegions: []string{randString(r, s)},
		}
	}

	randAWSProvider := func(r *rand.Rand, s int) *AWSProvider {
		if r.Float32() < 0.5 {
			accountID := json.Number(randString(r, s))
			if r.Float32() < 0.5 {
				return &AWSProvider{
					AccountID: &accountID,
					Region:    randStringPtr(r, s),
					Profile:   randStringPtr(r, s),
					Version:   randStringPtr(r, s),
				}
			} else {
				return &AWSProvider{
					AccountID: &accountID,
					Region:    randStringPtr(r, s),
					Role:      randStringPtr(r, s),
					Version:   randStringPtr(r, s),
				}
			}
		}
		return nil
	}

	randSnowflakeProvider := func(r *rand.Rand, s int) *SnowflakeProvider {
		if r.Float32() < 0.5 {
			return &SnowflakeProvider{
				Account: randStringPtr(r, size),
				Region:  randStringPtr(r, s),
				Role:    randStringPtr(r, s),
			}
		}
		return nil
	}

	randHerokuProvider := func(r *rand.Rand, s int) *HerokuProvider {
		if r.Float32() < 0.5 {
			return &HerokuProvider{}
		}
		return nil
	}

	randDatadogProvider := func(r *rand.Rand, s int) *DatadogProvider {
		if r.Float32() < 0.5 {
			return &DatadogProvider{}
		}
		return nil
	}

	randCommon := func(r *rand.Rand, s int) Common {
		var backendRole, backendProfile *string

		if r.Float32() < 0.5 {
			backendRole = randStringPtr(r, s)
		}

		if r.Float32() < 0.5 {
			backendProfile = randStringPtr(r, s)
		}

		c := Common{
			Backend: &Backend{
				Bucket:  randStringPtr(r, s),
				Region:  randStringPtr(r, s),
				Role:    backendRole,
				Profile: backendProfile,
			},
			ExtraVars: randStringMap(r, s),
			Owner:     randStringPtr(r, s),
			Project:   randStringPtr(r, s),
			Providers: &Providers{
				AWS:       randAWSProvider(r, s),
				Snowflake: randSnowflakeProvider(r, s),
				Okta:      randOktaProvider(r, s),
				Bless:     randBlessProvider(r, s),
				Heroku:    randHerokuProvider(r, s),
				Datadog:   randDatadogProvider(r, s),
			},
			TerraformVersion: randStringPtr(r, s),
		}

		if r.Float32() < 0.5 {
			c.Tools = &Tools{}
			if r.Float32() < 0.5 {
				c.Tools.TravisCI = &TravisCI{
					CommonCI: CommonCI{
						Enabled:     randBoolPtr(r, s),
						TestBuckets: randIntPtr(r, s),
					},
				}
			}
			if r.Float32() < 0.5 {
				p := r.Float32() < 0.5
				c.Tools.TfLint = &TfLint{
					Enabled: &p,
				}
			}
		}

		return c
	}

	conf.Version = 2

	conf.Defaults = Defaults{
		Common: randCommon(r, size),
	}

	// tools

	conf.Accounts = map[string]Account{}
	acctN := r.Intn(size)

	for i := 0; i < acctN; i++ {
		acctName := randString(r, size)
		conf.Accounts[acctName] = Account{
			Common: randCommon(r, size),
		}

	}

	conf.Envs = map[string]Env{}
	envN := r.Intn(size)

	for i := 0; i < envN; i++ {
		envName := randString(r, size)
		e := Env{
			Common: randCommon(r, size),
		}
		e.Components = map[string]Component{}
		compN := r.Intn(size)

		for i := 0; i < compN; i++ {
			compName := randString(r, size)
			e.Components[compName] = Component{
				Common: randCommon(r, size),
			}
		}
		conf.Envs[envName] = e

	}

	return reflect.ValueOf(conf)
}
