package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plugins"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v3"
)

// ReadConfig take a byte array as input and outputs a json or yaml config struct
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
	ComponentTemplates map[string]any     `yaml:"component_templates,omitempty"`
	Accounts           map[string]Account `yaml:"accounts,omitempty"`
	Defaults           Defaults           `yaml:"defaults" validate:"required"`
	Envs               map[string]Env     `yaml:"envs,omitempty"`
	Global             Component          `yaml:"global,omitempty"`
	Modules            map[string]Module  `yaml:"modules,omitempty"`
	Plugins            Plugins            `yaml:"plugins,omitempty"`
	Version            int                `validate:"required,eq=2"`
	TFE                *TFE               `yaml:"tfe,omitempty"`
}

type TFE struct {
	Component                      `yaml:",inline"`
	ReadTeams                      *[]string `yaml:"read_teams,omitempty"`
	Branch                         *string   `yaml:"branch,omitempty"`
	GithubOrg                      *string   `yaml:"gh_org,omitempty"`
	GithubRepo                     *string   `yaml:"gh_repo,omitempty"`
	TFEOrg                         string    `yaml:"tfe_org,omitempty"`
	SSHKeyName                     *string   `yaml:"ssh_key_name,omitempty"`
	ExcludedGithubRequiredChecks   *[]string `yaml:"excluded_gh_required_checks,omitempty"`
	AdditionalGithubRequiredChecks *[]string `yaml:"additional_gh_required_checks,omitempty"`
}

type ExtraTemplate struct {
	Overwrite *bool
	Content   *string
}

type Common struct {
	Backend                  *Backend                  `yaml:"backend,omitempty"`
	ExtraVars                map[string]string         `yaml:"extra_vars,omitempty"`
	Owner                    *string                   `yaml:"owner,omitempty"`
	Project                  *string                   `yaml:"project,omitempty"`
	Providers                *Providers                `yaml:"providers,omitempty"`
	DependsOn                *DependsOn                `yaml:"depends_on,omitempty"`
	TerraformVersion         *string                   `yaml:"terraform_version,omitempty"`
	Tools                    *Tools                    `yaml:"tools,omitempty"`
	NeedsAWSAccountsVariable *bool                     `yaml:"needs_aws_accounts_variable,omitempty"`
	ExtraTemplates           *map[string]ExtraTemplate `yaml:"extra_templates,omitempty"`
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

	SSHKeySecrets []string  `yaml:"ssh_key_secrets"`
	RunsOn        *[]string `yaml:"runs_on,omitempty"`
}

type Env struct {
	Common `yaml:",inline"`

	Components map[string]Component `yaml:"components,omitempty"`
}

type Component struct {
	Common `yaml:",inline"`

	EKS             *EKSConfig         `yaml:"eks,omitempty"`
	Kind            *ComponentKind     `yaml:"kind,omitempty"`
	ModuleSource    *string            `yaml:"module_source,omitempty"`
	ModuleName      *string            `yaml:"module_name,omitempty"`
	ProviderAliases *map[string]string `yaml:"provider_aliases,omitempty"`
	Variables       *[]string          `yaml:"variables,omitempty"`
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
	Helm       *HelmProvider       `yaml:"helm,omitempty"`
	Okta       *OktaProvider       `yaml:"okta,omitempty"`
	OpsGenie   *OpsGenieProvider   `yaml:"opsgenie,omitempty"`
	Pagerduty  *PagerdutyProvider  `yaml:"pagerduty,omitempty"`
	Sentry     *SentryProvider     `yaml:"sentry,omitempty"`
	Snowflake  *SnowflakeProvider  `yaml:"snowflake,omitempty"`
	Tfe        *TfeProvider        `yaml:"tfe,omitempty"`
	Kubectl    *KubectlProvider    `yaml:"kubectl,omitempty"`
}

type AssertProvider struct {
	CommonProvider `yaml:",inline"`
}

// CommonProvider encapsulates common properties across providers
// TODO refactor other providers to use CommonProvider inline
type CommonProvider struct {
	CustomProvider *bool   `yaml:"custom_provider,omitempty"`
	Enabled        *bool   `yaml:"enabled,omitempty"`
	Version        *string `yaml:"version,omitempty"`
}

// Auth0Provider is the terraform provider for the Auth0 service.
type Auth0Provider struct {
	CommonProvider `yaml:",inline"`
	Domain         *string `yaml:"domain,omitempty"`
	Source         *string `yaml:"source,omitempty"`
}

// OktaProvider is an okta provider
type OktaProvider struct {
	CommonProvider `yaml:",inline"`
	// the okta provider is optional (above) but if supplied you must set an OrgName

	RegistryNamespace *string `yaml:"registry_namespace"` // for forked providers

	OrgName *string `yaml:"org_name,omitempty"`
	BaseURL *string `yaml:"base_url,omitempty"`
}

// BlessProvider allows for terraform-provider-bless configuration
type BlessProvider struct {
	CommonProvider `yaml:",inline"`
	// the bless provider is optional (above) but if supplied you must set a region and aws_profile
	AdditionalRegions []string `yaml:"additional_regions,omitempty"`
	AWSProfile        *string  `yaml:"aws_profile,omitempty"`
	AWSRegion         *string  `yaml:"aws_region,omitempty"`
	RoleArn           *string  `yaml:"role_arn,omitempty"`
}

type AWSProvider struct {
	CommonProvider `yaml:",inline"`
	// the aws provider is optional (above) but if supplied you must set account id and region
	AccountID         *json.Number `yaml:"account_id,omitempty"`
	AdditionalRegions []string     `yaml:"additional_regions,omitempty"`
	Profile           *string      `yaml:"profile,omitempty"`
	Region            *string      `yaml:"region,omitempty"`
	Role              *string      `yaml:"role,omitempty"` // FIXME validate format
	// HACK HACK(el): we can configure additional, aliased, AWS providers for other accounts
	// 								A map of alias_name to provider configuration
	AdditionalProviders map[string]*AWSProvider `yaml:"additional_providers,omitempty"`
}

type GithubProvider struct {
	CommonProvider `yaml:",inline"`
	Organization   *string `yaml:"organization,omitempty"`
	BaseURL        *string `yaml:"base_url,omitempty"`
}

type GrafanaProvider struct {
	CommonProvider `yaml:",inline"`
}

type SnowflakeProvider struct {
	CommonProvider `yaml:",inline"`
	Account        *string `yaml:"account,omitempty"`
	Role           *string `yaml:"role,omitempty"`
	Region         *string `yaml:"region,omitempty"`
}

type HerokuProvider struct {
	CommonProvider `yaml:",inline"`
}

type DatadogProvider struct {
	CommonProvider `yaml:",inline"`
}

type PagerdutyProvider struct {
	CommonProvider `yaml:",inline"`
}

type OpsGenieProvider struct {
	CommonProvider `yaml:",inline"`
}

type DatabricksProvider struct {
	CommonProvider `yaml:",inline"`
}

type SentryProvider struct {
	CommonProvider `yaml:",inline"`

	BaseURL *string `yaml:"base_url,omitempty"`
}

type TfeProvider struct {
	CommonProvider `yaml:",inline"`

	Hostname *string `yaml:"hostname,omitempty"`
}

type KubernetesProvider struct {
	CommonProvider       `yaml:",inline"`
	ClusterComponentName *string `yaml:"cluster_component_name,omitempty"`
}

type HelmProvider struct {
	CommonProvider `yaml:",inline"`
}

type KubectlProvider struct {
	CommonProvider `yaml:",inline"`
}

// Backend is used to configure a terraform backend
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
