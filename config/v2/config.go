package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plugins"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v3"
)

//ReadConfig take a byte array as input and outputs a json or yaml config struct
func ReadConfig(fs afero.Fs, b []byte, configFile string) (*Config, error) {
	var e error
	c := &Config{}

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

func (c *Config) Write(fs afero.Fs, path string) error {
	yamlConfigFile, err := fs.Create("fogg.yml")
	if err != nil {
		return errs.WrapInternal(err, "unable to create config file fogg.yml")
	}
	defer yamlConfigFile.Close()

	encoder := yaml.NewEncoder(yamlConfigFile)
	encoder.SetIndent(2)

	return encoder.Encode(c)
}

type Config struct {
	Accounts map[string]Account `yaml:"accounts,omitempty"`
	Defaults Defaults           `yaml:"defaults" validate:"required"`
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
	DependsOn        *DependsOn        `yaml:"depends_on,omitempty"`
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

	SSHKeyFingerprints []string `yaml:"ssh_key_fingerprints,omitempty"`
}

type GitHubActionsCI struct {
	CommonCI `yaml:",inline"`

	SSHKeySecrets []string `yaml:"ssh_key_secrets"`
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
	ModuleName   *string        `yaml:"module_name,omitempty"`
}

type Providers struct {
	Assert     *AssertProvider     `yaml:"assert,omitempty"`
	Auth0      *Auth0Provider      `yaml:"auth0,omitempty"`
	AWS        *AWSProvider        `yaml:"aws,omitempty"`
	Bless      *BlessProvider      `yaml:"bless,omitempty"`
	Databricks *DatabricksProvider `yaml:"databricks,omitempty"`
	Datadog    *DatadogProvider    `yaml:"datadog,omitempty"`
	Github     *GithubProvider     `yaml:"github,omitempty"`
	Grafana    *GrafanaProvider    `yaml:"grafana,omitempty"`
	Heroku     *HerokuProvider     `yaml:"heroku,omitempty"`
	Kubernetes *KubernetesProvider `yaml:"kubernetes,omitempty"`
	Okta       *OktaProvider       `yaml:"okta,omitempty"`
	OpsGenie   *OpsGenieProvider   `yaml:"opsgenie,omitempty"`
	Pagerduty  *PagerdutyProvider  `yaml:"pagerduty,omitempty"`
	Sentry     *SentryProvider     `yaml:"sentry,omitempty"`
	Snowflake  *SnowflakeProvider  `yaml:"snowflake,omitempty"`
	Tfe        *TfeProvider        `yaml:"tfe,omitempty"`
}

type AssertProvider struct {
	Version *string `yaml:"version,omitempty"`
}

// CommonProvider encapsulates common properties across providers
// TODO refactor other providers to use CommonProvider inline
type CommonProvider struct {
	Enabled *bool   `yaml:"enabled,omitempty"`
	Version *string `yaml:"version,omitempty"`
}

//Auth0Provider is the terraform provider for the Auth0 service.
type Auth0Provider struct {
	Version *string `yaml:"version,omitempty"`
	Domain  *string `yaml:"domain,omitempty"`
}

// OktaProvider is an okta provider
type OktaProvider struct {
	// the okta provider is optional (above) but if supplied you must set an OrgName

	// TODO refactor to get these from CommonProvider
	Version           *string `yaml:"version,omitempty"`
	RegistryNamespace *string `yaml:"registry_namespace"` // for forked providers

	OrgName *string `yaml:"org_name,omitempty"`
	BaseURL *string `yaml:"base_url,omitempty"`
}

// BlessProvider allows for terraform-provider-bless configuration
type BlessProvider struct {
	// the bless provider is optional (above) but if supplied you must set a region and aws_profile
	AdditionalRegions []string `yaml:"additional_regions,omitempty"`
	AWSProfile        *string  `yaml:"aws_profile,omitempty"`
	AWSRegion         *string  `yaml:"aws_region,omitempty"`
	RoleArn           *string  `yaml:"role_arn,omitempty"`
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
	// HACK HACK(el): we can configure additional, aliased, AWS providers for other accounts
	// 								A map of alias_name to provider configuration
	AdditionalProviders map[string]*AWSProvider `yaml:"additional_providers,omitempty"`
}

type GithubProvider struct {
	Organization *string `yaml:"organization,omitempty"`
	BaseURL      *string `yaml:"base_url,omitempty"`
	Version      *string `yaml:"version,omitempty"`
}

type GrafanaProvider struct {
	CommonProvider `yaml:",inline"`
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

type PagerdutyProvider struct {
	Version *string `yaml:"version,omitempty"`
}

type OpsGenieProvider struct {
	Version *string `yaml:"version,omitempty"`
}

type DatabricksProvider struct {
	Version *string `yaml:"version,omitempty"`
}

type SentryProvider struct {
	Version *string `yaml:"version,omitempty"`
	BaseURL *string `yaml:"base_url,omitempty"`
}

type TfeProvider struct {
	CommonProvider `yaml:",inline"`

	Hostname *string `yaml:"hostname,omitempty"`
}

type KubernetesProvider struct {
	CommonProvider `yaml:",inline"`
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
	Env            map[string]string           `yaml:"env,omitempty"`
}

type DependsOn struct {
	Accounts   []string `yaml:"accounts"`
	Components []string `yaml:"components"`
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

type componentInfo struct {
	Kind string
	Name string
	Env  string
}

// PathToComponentType given a path, return information about that component
func (c Config) PathToComponentType(path string) (componentInfo, error) { //nolint
	t := componentInfo{}

	path = strings.TrimRight(path, "/")
	pathParts := strings.Split(path, "/")
	switch len(pathParts) {
	case 3:
		accountName := pathParts[2]
		if _, found := c.Accounts[accountName]; !found {
			return t, fmt.Errorf("could not find account %s", accountName)
		}
		t.Kind = "accounts"
		t.Name = accountName

		return t, nil
	case 4:
		envName := pathParts[2]
		componentName := pathParts[3]

		env, envFound := c.Envs[envName]

		if !envFound {
			return t, fmt.Errorf("could not find env %s", envName)
		}

		_, componentFound := env.Components[componentName]

		if !componentFound {
			return t, fmt.Errorf("could not find component %s in env %s", componentName, envName)
		}
		t.Kind = "envs"
		t.Name = componentName
		t.Env = envName
		return t, nil
	default:
		return t, fmt.Errorf("could not figure out component for path %s", path)
	}
}

const (
	// DefaultComponentKind defaults to terraform component
	DefaultComponentKind ComponentKind = "terraform"
	// ComponentKindTerraform is a terraform component
	ComponentKindTerraform = DefaultComponentKind
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

	randBoolPtr := func(r *rand.Rand) *bool {
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

	randAuth0Provider := func(r *rand.Rand, s int) *Auth0Provider {
		return &Auth0Provider{
			Version: randStringPtr(r, s),
			Domain:  randStringPtr(r, s),
		}
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
			}
			return &AWSProvider{
				AccountID: &accountID,
				Region:    randStringPtr(r, s),
				Role:      randStringPtr(r, s),
				Version:   randStringPtr(r, s),
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

	randHerokuProvider := func(r *rand.Rand) *HerokuProvider {
		if r.Float32() < 0.5 {
			return &HerokuProvider{}
		}
		return nil
	}

	randDatadogProvider := func(r *rand.Rand) *DatadogProvider {
		if r.Float32() < 0.5 {
			return &DatadogProvider{}
		}
		return nil
	}

	randPagerdutyProvider := func(r *rand.Rand) *PagerdutyProvider {
		if r.Float32() < 0.5 {
			return &PagerdutyProvider{}
		}
		return nil
	}

	randOpsGenieProvider := func(r *rand.Rand) *OpsGenieProvider {
		if r.Float32() < 0.5 {
			return &OpsGenieProvider{}
		}
		return nil
	}

	randDatabricksProvider := func(r *rand.Rand) *DatabricksProvider {
		if r.Float32() < 0.5 {
			return &DatabricksProvider{}
		}
		return nil
	}

	randKubernetesProvider := func(r *rand.Rand) *KubernetesProvider {
		if r.Float32() < 0.5 {
			return &KubernetesProvider{}
		}
		return nil
	}

	randGrafanaProvider := func(r *rand.Rand) *GrafanaProvider {
		if r.Float32() < 0.5 {
			return &GrafanaProvider{}
		}
		return nil
	}

	randSentryProvider := func(r *rand.Rand) *SentryProvider {
		if r.Float32() < 0.5 {
			return &SentryProvider{}
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
				Auth0:      randAuth0Provider(r, s),
				AWS:        randAWSProvider(r, s),
				Bless:      randBlessProvider(r, s),
				Datadog:    randDatadogProvider(r),
				Pagerduty:  randPagerdutyProvider(r),
				Databricks: randDatabricksProvider(r),
				Grafana:    randGrafanaProvider(r),
				Heroku:     randHerokuProvider(r),
				Kubernetes: randKubernetesProvider(r),
				Okta:       randOktaProvider(r, s),
				OpsGenie:   randOpsGenieProvider(r),
				Sentry:     randSentryProvider(r),
				Snowflake:  randSnowflakeProvider(r, s),
			},
			TerraformVersion: randStringPtr(r, s),
		}

		if r.Float32() < 0.5 {
			c.Tools = &Tools{}
			if r.Float32() < 0.5 {
				c.Tools.TravisCI = &TravisCI{
					CommonCI: CommonCI{
						Enabled:     randBoolPtr(r),
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
