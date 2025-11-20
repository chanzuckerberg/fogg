package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plugins"
	"github.com/runatlantis/atlantis/server/core/config/raw"
	"github.com/sirupsen/logrus"
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
		if e == nil && c.ConfDir != nil && *c.ConfDir != "" {
			logrus.Debugf("Conf dir is %q\n", *c.ConfDir)
			e = ReadConfDir(fs, c)
		}

	default:
		return nil, errs.NewUserf("File type %s is not supported", ext)
	}
	return c, e
}

func ReadConfDir(fs afero.Fs, c *Config) error {
	info, e := fs.Stat(*c.ConfDir)
	if e != nil {
		return errs.WrapUserf(e, "unable to find conf_dir %q", *c.ConfDir)
	}
	if !info.IsDir() {
		return errs.WrapUserf(e, "conf_dir %q must be a directory", *c.ConfDir)
	}
	logrus.Debugf("Walking Conf dir %q\n", *c.ConfDir)
	partialConfigs := []*Config{c}
	e = afero.Walk(fs, *c.ConfDir, func(path string, info os.FileInfo, err error) error {
		// TODO: ignore more files?
		if info.IsDir() {
			logrus.Debugf("Ignoring %q\n", path)
			return nil
		}
		logrus.Debugf("Opening %q\n", path)
		partial, e := fs.Open(path)
		if e != nil {
			logrus.Debugf("Ignoring error opening %q\n", path)
			return nil
		}
		b, e := io.ReadAll(partial)
		if e != nil {
			return errs.WrapUserf(e, "unable to read partial config %q", path)
		}
		pc, e := ReadConfig(fs, b, path)
		if e != nil {
			return errs.WrapUserf(e, "unable to parse partial config %q", path)
		}
		logrus.Debugf("appending partialConfig %q\n", path)
		partialConfigs = append(partialConfigs, pc)
		return nil
	})
	if e != nil {
		return errs.WrapUserf(e, "unable to walk conf_dir %q", *c.ConfDir)
	}
	// merge partialConfigs into c
	mergeConfigs(partialConfigs...)
	return e
}

func mergeConfigs(confs ...*Config) {
	if len(confs) < 2 {
		return
	}
	mergedConfig, tail := confs[0], confs[1:]
	for _, pc := range tail {
		if mergedConfig.Accounts == nil {
			if pc.Accounts != nil {
				mergedConfig.Accounts = pc.Accounts
			}
		} else {
			if pc.Accounts != nil {
				maps.Copy(mergedConfig.Accounts, pc.Accounts)
			}
		}

		if mergedConfig.Envs == nil {
			if pc.Envs != nil {
				mergedConfig.Envs = pc.Envs
			}
		} else {
			if pc.Envs != nil {
				maps.Copy(mergedConfig.Envs, pc.Envs)
			}
		}

		if mergedConfig.Modules == nil {
			if pc.Modules != nil {
				mergedConfig.Modules = pc.Modules
			}
		} else {
			if pc.Modules != nil {
				maps.Copy(mergedConfig.Modules, pc.Modules)
			}
		}
	}
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
	Modules  map[string]Module  `yaml:"modules,omitempty"` // BUG: order is important
	Plugins  Plugins            `yaml:"plugins,omitempty"`
	Version  int                `validate:"required,eq=2"`
	TFE      *TFE               `yaml:"tfe,omitempty"`
	ConfDir  *string            `yaml:"conf_dir,omitempty"`
	Turbo    *TurboConfig       `yaml:"turbo,omitempty"`
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

type Common struct {
	Backend           *Backend                    `yaml:"backend,omitempty"`
	ExtraVars         map[string]string           `yaml:"extra_vars,omitempty"`
	Owner             *string                     `yaml:"owner,omitempty"`
	Project           *string                     `yaml:"project,omitempty"`
	Providers         *Providers                  `yaml:"providers,omitempty"`
	RequiredProviders map[string]*GenericProvider `yaml:"required_providers,omitempty"`
	DependsOn         *DependsOn                  `yaml:"depends_on,omitempty"`
	TerraformVersion  *string                     `yaml:"terraform_version,omitempty"`
	Tools             *Tools                      `yaml:"tools,omitempty"`
	// Store output for Integrations (only ssm supported atm)
	IntegrationRegistry *string     `yaml:"integration_registry,omitempty"`
	Grid                *GridConfig `yaml:"grid,omitempty"`
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
	Atlantis        *Atlantis        `yaml:"atlantis,omitempty"`
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

type Atlantis struct {
	// enable Atlantis integration
	// default: false
	Enabled *bool `yaml:"enabled,omitempty"`
	// list of module source prefixes for auto plan when modified
	// default: "terraform/modules/"
	//
	// Deprecated: autoplan now auto detects if module is local
	ModulePrefixes []string `yaml:"module_prefixes,omitempty"`
	// autoplan remote-states (only if depends_on is provided)
	// default: false
	AutoplanRemoteStates *bool `yaml:"autoplan_remote_states,omitempty"`
	// Raw atlantis RepoCfg struct
	raw.RepoCfg `yaml:",inline"`
}

type Env struct {
	Common `yaml:",inline"`

	NoGlobal   bool                 `yaml:"no_global,omitempty"`
	Components map[string]Component `yaml:"components,omitempty"`
}

// TODO: Support cdktf depedencies from a private registry/scope
type Component struct {
	Common `yaml:",inline"`

	EKS                  *EKSConfig             `yaml:"eks,omitempty"`
	Kind                 *ComponentKind         `yaml:"kind,omitempty"`
	ModuleSource         *string                `yaml:"module_source,omitempty"`
	ModuleName           *string                `yaml:"module_name,omitempty"`
	ModuleForEach        *string                `yaml:"module_for_each,omitempty"`
	ProvidersMap         map[string]string      `yaml:"module_providers,omitempty"`
	Variables            []string               `yaml:"variables,omitempty"`
	Outputs              []string               `yaml:"outputs,omitempty"`
	Modules              []ComponentModule      `yaml:"modules,omitempty"`
	CdktfDependencies    []JavascriptDependency `yaml:"cdktf_dependencies,omitempty"`     // Optional additional component dev dependencies, default: []
	CdktfDevDependencies []JavascriptDependency `yaml:"cdktf_dev_dependencies,omitempty"` // Optional additional component dev dependencies, default: []
	PackageJsonFields    map[string]any         `yaml:"package_json,omitempty"`           // Fields to add into package.json, only used by non-cdktf components
}

type GridConfig struct {
	Enabled  *bool   `yaml:"enabled,omitempty"`
	GUID     *string `yaml:"guid,omitempty"`
	Endpoint *string `yaml:"endpoint,omitempty"`
}

type ComponentModule struct {
	// Source for Terraform module as supported by Terraform
	Source *string `yaml:"source,omitempty"`
	// Version for Terraform module as supported by Terraform
	Version *string `yaml:"version,omitempty"`
	// Name for generated module block, defaults to Source stripped from special characters
	Name *string `yaml:"name,omitempty"`
	// Prefix for all generated input and output placeholder to handle overlapping references
	Prefix *string `yaml:"prefix,omitempty"`
	// Variables to limit generated input placeholders (and use module defaults for others)
	Variables []string `yaml:"variables,omitempty"`
	// Outputs list to limit generated component outputs
	Outputs []string `yaml:"outputs,omitempty"`
	// Integration Registry config
	Integration *ModuleIntegrationConfig `yaml:"integration,omitempty"`
	// Optional mapping of providers https://developer.hashicorp.com/terraform/language/meta-arguments/module-providers
	ProvidersMap map[string]string `yaml:"providers,omitempty"`
	// For Each metadata argument https://developer.hashicorp.com/terraform/language/modules/syntax#meta-arguments
	ForEach *string `yaml:"for_each,omitempty"`
	// Dependencies of this module
	DependsOn []string `yaml:"depends_on,omitempty"`
}

type ModuleIntegrationConfig struct {
	// Mode only "none" | "selected" | "all" supported
	// default = "none", anything else is treated as "all"
	Mode *string `yaml:"mode,omitempty"`
	// A default golang format string for output integration
	// omitted format is "module.module_name.output_name"
	Format *string `yaml:"format,omitempty"`
	// Drop prefix used for input and output placeholders from parameter path
	DropPrefix bool `yaml:"drop_prefix,omitempty"`
	// Drop component from parameter path (only uses env)
	DropComponent bool `yaml:"drop_component,omitempty"`
	// Infix path for all outputs
	PathInfix *string `yaml:"path_infix,omitempty"`
	// Resource providers to publish outputs via https://developer.hashicorp.com/terraform/language/meta-arguments/resource-provider
	Providers []string `yaml:"providers,omitempty"`
	// Map for outputs into Integration Registry
	OutputsMap map[string]*IntegrationRegistryMap `yaml:"outputs_map,omitempty"`
}

type IntegrationRegistryMap struct {
	// A golang format string
	Format *string `yaml:"format,omitempty"`
	// Drop component from parameter path (only use env)
	DropComponent *bool `yaml:"drop_component,omitempty"`
	// Path to store outputs under
	Path *string `yaml:"path,omitempty"`
	// Add for each configuration
	ForEach bool `yaml:"for_each,omitempty"`
	// Create for_each outputs under this path
	PathForEach *string `yaml:"path_for_each,omitempty"`
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
	Sops       *SopsProvider       `yaml:"sops,omitempty"`
}

type AssertProvider struct {
	CommonProvider `yaml:",inline"`
}

// CommonProvider encapsulates common properties across providers
type CommonProvider struct {
	CustomProvider *bool   `yaml:"custom_provider,omitempty" json:"custom_provider,omitempty"`
	Enabled        *bool   `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Version        *string `yaml:"version,omitempty" json:"version,omitempty"`
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
	AccountID         *json.Number            `yaml:"account_id,omitempty"`
	AdditionalRegions []string                `yaml:"additional_regions,omitempty"`
	Profile           *string                 `yaml:"profile,omitempty"`
	Region            *string                 `yaml:"region,omitempty"`
	Role              *string                 `yaml:"role,omitempty"` // FIXME validate format
	DefaultTags       *AWSProviderDefaultTags `yaml:"default_tags,omitempty"`
	IgnoreTags        *AWSProviderIgnoreTags  `yaml:"ignore_tags,omitempty"`
	// HACK HACK(el): we can configure additional, aliased, AWS providers for other accounts
	// 								A map of alias_name to provider configuration
	AdditionalProviders map[string]*AWSProvider `yaml:"additional_providers,omitempty"`
}

type AWSProviderDefaultTags struct {
	Enabled *bool `yaml:"enabled,omitempty"`
	// List of exact resource tag keys to ignore across all resources handled by this provider.
	Tags map[string]string `yaml:"tags,omitempty"`
}

type AWSProviderIgnoreTags struct {
	Enabled *bool `yaml:"enabled,omitempty"`
	// List of exact resource tag keys to ignore across all resources handled by this provider.
	Keys        []string `yaml:"keys,omitempty"`
	KeyPrefixes []string `yaml:"key_prefixes,omitempty"`
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

type SopsProvider struct {
	CommonProvider `yaml:",inline"`
}

type KubernetesProvider struct {
	CommonProvider `yaml:",inline"`
}

// GenericProvider is a generic terraform provider.
type GenericProvider struct {
	CommonProvider `yaml:",inline" json:",inline"`
	Source         string         `yaml:"source" json:"source"`
	Config         map[string]any `yaml:"config" json:"config"`
}

// Backend is used to configure a terraform backend
type Backend struct {
	Kind *string `yaml:"kind,omitempty" validate:"omitempty,oneof=s3 remote http"`

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

	// fields used for http backend
	BaseAddress   *string `yaml:"base_address,omitempty"`
	Address       *string `yaml:"address,omitempty"`
	LockAddress   *string `yaml:"lock_address,omitempty"`
	UnlockAddress *string `yaml:"unlock_address,omitempty"`
	UpdateMethod  *string `yaml:"update_method,omitempty"`
	LockMethod    *string `yaml:"lock_method,omitempty"`
	UnlockMethod  *string `yaml:"unlock_method,omitempty"`
	Username      *string `yaml:"username,omitempty"`
	Password      *string `yaml:"password,omitempty"`
}

// Module is a module
type Module struct {
	Kind                  *ModuleKind            `yaml:"kind,omitempty"` // terraform or cdktf
	TerraformVersion      *string                `yaml:"terraform_version,omitempty"`
	PackageName           *string                `yaml:"package_name,omitempty"`
	Publish               *bool                  `yaml:"publish,omitempty"`
	Owner                 *string                `yaml:"author,omitempty"`
	CdktfDependencies     []JavascriptDependency `yaml:"cdktf_dependencies,omitempty"`      // Optional additional module dependencies, default: []
	CdktfDevDependencies  []JavascriptDependency `yaml:"cdktf_dev_dependencies,omitempty"`  // Optional additional module dev dependencies, default: []
	CdktfPeerDependencies []JavascriptDependency `yaml:"cdktf_peer_dependencies,omitempty"` // Optional additional module peer dependencies, default: []
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
	AWSRegion      *string                     `yaml:"aws_region,omitempty"`
	TestBuckets    *int                        `yaml:"test_buckets,omitempty"`
	Command        *string                     `yaml:"command,omitempty"`
	Buildevents    *bool                       `yaml:"buildevents,omitempty"`
	Providers      map[string]CIProviderConfig `yaml:"providers,omitempty"`
	Env            map[string]string           `yaml:"env,omitempty"`
	PreCommit      *PreCommitSetup             `yaml:"pre_commit,omitempty"`
}

type PreCommitSetup struct {
	Enabled bool    `yaml:"enabled"`
	Version *string `yaml:"version,omitempty"`
	// Pip requirements for better CI cache usage
	PipCache map[string]string `yaml:"pip_cache,omitempty"`
	// Additional CI steps for pre-commit setup
	GitHubActionSteps []GitHubActionStep `yaml:"github_actions_setup,omitempty"`
	// Extra Args to pass into pre-commit call
	ExtraArgs []string `yaml:"extra_args,omitempty"`
	// Simplified pre-commit config (with custom fields such as `skip_in_make`)
	Config *PreCommitConfig `yaml:"config,omitempty"`
}

type GitHubActionStep struct {
	Name *string           `yaml:"name,omitempty"`
	Uses *string           `yaml:"uses,omitempty"`
	With map[string]string `yaml:"with,omitempty"`
	Run  *string           `yaml:"run,omitempty"`
}

type PreCommitConfig struct {
	Files    string `yaml:"files,omitempty"`
	Exclude  string `yaml:"exclude,omitempty"`
	FailFast *bool  `yaml:"fail_fast,omitempty"`
	Repos    []Repo `yaml:"repos"`
}

type Repo struct {
	Repo  string  `yaml:"repo"`
	Rev   *string `yaml:"rev,omitempty"`
	Hooks []Hook  `yaml:"hooks"`
}

type Hook struct {
	ID                     string   `yaml:"id"`
	Name                   *string  `yaml:"name,omitempty"`
	Alias                  *string  `yaml:"alias,omitempty"`
	Args                   []string `yaml:"args,omitempty"`
	Exclude                *string  `yaml:"exclude,omitempty"`
	Files                  *string  `yaml:"files,omitempty"`
	AdditionalDependencies []string `yaml:"additional_dependencies,omitempty"`
	RequireSerial          *bool    `yaml:"require_serial,omitempty"`
	// required for repo = "local" hooks
	Entry    *string `yaml:"entry,omitempty"`
	Language *string `yaml:"language,omitempty"`
	// skip in make target (for quick pre-commit autofix after fogg apply)
	// this field is set to nil to avoid pre-commit warnings on invalid field
	// see:plan/ci.go -> buildGithubActionsPreCommitConfig
	SkipInMake *bool `yaml:"skip_in_make,omitempty"`
}

type DependsOn struct {
	Accounts   DependencyList `yaml:"accounts"`
	Components DependencyList `yaml:"components"`
	//RelativeGlobs to the component
	RelativeGlobs []string `yaml:"relative_globs"`
	//Absolute file paths,
	//fogg validates their existence and
	//generates locals block for their content
	Files []string `yaml:"files"`
}

// DependencyList is a map of dependencies to their outputs
type DependencyList map[string][]string

func (d *DependencyList) UnmarshalYAML(value *yaml.Node) error {
	var list []string
	if err := value.Decode(&list); err == nil {
		*d = make(map[string][]string)
		for _, item := range list {
			(*d)[item] = []string{}
		}
		return nil
	}

	var m map[string][]string
	if err := value.Decode(&m); err == nil {
		*d = m
		return nil
	}

	return fmt.Errorf("failed to unmarshal DependsOn: expected list of strings or map of string to list of strings")
}

type CIProviderConfig struct {
	Disabled bool `yaml:"disabled,omitempty"`
}

type TfLint struct {
	Enabled *bool `yaml:"enabled,omitempty"`
}

type TurboConfig struct {
	Enabled  *bool   `yaml:"enabled,omitempty"`   // Enable Turbo, default: false
	Version  *string `yaml:"version,omitempty"`   // Optional Turbo version, default: "^2.0.6"
	RootName *string `yaml:"root_name,omitempty"` // Optional Name for the root package, default: "fogg-monorepo"
	SCMBase  *string `yaml:"scm_base,omitempty"`  // Optional Git comparison base override, default: "main"

	Scopes          []JavascriptPackageScope `yaml:"scopes,omitempty"`           // Optional additional scopes, default: []
	DevDependencies []JavascriptDependency   `yaml:"dev_dependencies,omitempty"` // Optional additional root dev dependencies, default: []
}

type JavascriptDependency struct {
	Name    string `yaml:"name"`    // npm package name
	Version string `yaml:"version"` // npm package version
}

type JavascriptPackageScope struct {
	Name        string `yaml:"name"`         // name for example "@vincenthsh"
	RegistryUrl string `yaml:"registry_url"` // registry url for example "https://npm.pkg.github.com"
	AlwaysAuth  bool   `yaml:"always_auth"`  // always auth flag, default: false
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

// ComponentKind is the kind of this component
type ModuleKind string

// GetOrDefault gets the component kind or default
func (ck *ModuleKind) GetOrDefault() ModuleKind {
	if ck == nil || *ck == "" {
		return DefaultModuleKind
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
	// ComponentKindCDKTF is a CDKTF component
	ComponentKindCDKTF ComponentKind = "cdktf"
	// ComponentKindTerraConstruct is a CDKTF component using the TerraConstructs framework
	ComponentKindTerraConstruct ComponentKind = "terraconstruct"
	// DefaultComponentKind defaults to terraform component
	DefaultModuleKind ModuleKind = "terraform"
	// ModuleKindTerraform is a terraform Module
	ModuleKindTerraform = DefaultModuleKind
	// ModuleKindCDKTF is a CDKTF Module
	ModuleKindCDKTF ModuleKind = "cdktf"
)
