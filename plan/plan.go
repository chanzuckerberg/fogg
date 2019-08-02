package plan

import (
	"encoding/json"
	"errors"
	"fmt"

	v1 "github.com/chanzuckerberg/fogg/config/v1"
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/util"
	"gopkg.in/yaml.v3"
)

// Plan represents a set of actions to take
type Plan struct {
	Accounts map[string]Account `json:"account" yaml:"account"`
	Atlantis Atlantis           `json:"atlantis" yaml:"atlantis"`
	Envs     map[string]Env     `json:"envs" yaml:"envs"`
	Global   Component          `json:"global" yaml:"global"`
	Modules  map[string]Module  `json:"modules" yaml:"modules"`
	TravisCI TravisCI           `json:"travis_ci" yaml:"travis_ci"`
	Version  string             `json:"version" yaml:"version"`
}

// Common represents common fields
type Common struct {
	PathToRepoRoot   string `yaml:"path_to_repo_root"`
	TerraformVersion string `yaml:"terraform_version"`
}

//ComponentCommon represents common fields for components
type ComponentCommon struct {
	Common `yaml:",inline"`

	Atlantis  AtlantisComponent `yaml:"atlantis"`
	Backend   AWSBackend        `json:"backend" yaml:"backend"`
	ExtraVars map[string]string `json:"extra_vars" yaml:"extra_vars"`
	Owner     string            `json:"owner" yaml:"owner"`
	Project   string            `json:"project" yaml:"project"`
	Providers Providers         `json:"providers" yaml:"providers"`
	TfLint    TfLint            `json:"tf_lint" yaml:"tf_lint"`
	TravisCI  TravisComponent
}

type AtlantisComponent struct {
	Enabled  bool   `yaml:"enabled"`
	RoleName string `yaml:"role_name"`
	RolePath string `yaml:"role_path"`
}

type TravisComponent struct {
	Enabled        bool
	AWSProfileName string
	AWSRoleName    string
	AWSAccountID   string
	Command        string
}

type Providers struct {
	AWS       *AWSProvider       `yaml:"aws"`
	Github *GithubProvider `yaml:"github"`
	Snowflake *SnowflakeProvider `yaml:"snowflake"`
	Bless     *BlessProvider     `yaml:"bless"`
	Okta      *OktaProvider      `yaml:"okta"`
}

//AWSProvider represents AWS provider configuration
type AWSProvider struct {
	AccountID         json.Number `yaml:"account_id"`
	Profile           string      `yaml:"profile"`
	Version           string      `yaml:"version"`
	Region            string      `yaml:"region"`
	AdditionalRegions []string    `yaml:"additional_regions"`
}

// GithubProvider represents a configuration of a github provider
type GithubProvider struct {
	Organization *string `yaml:"organization"`
	BaseURL      *string `yaml:"base_url"`
}

//SnowflakeProvider represents Snowflake DB provider configuration
type SnowflakeProvider struct {
	Account string  `json:"account,omitempty" yaml:"account,omitempty"`
	Role    string  `json:"role,omitempty" yaml:"role,omitempty"`
	Region  string  `json:"region,omitempty" yaml:"region,omitempty"`
	Version *string `json:"version,omitempty" yaml:"version,omitempty"`
}

//OktaProvider represents Okta configuration
type OktaProvider struct {
	OrgName string  `json:"org_name,omitempty"`
	Version *string `json:"version,omitempty"`
}

//BlessProvider represents Bless ssh provider configuration
type BlessProvider struct {
	AdditionalRegions []string `json:"additional_regions,omitempty" yaml:"additional_regions,omitempty"`
	AWSProfile        string   `json:"aws_profile,omitempty" yaml:"aws_profile,omitempty"`
	AWSRegion         string   `json:"aws_region,omitempty" yaml:"aws_region,omitempty"`
	Version           *string  `json:"version,omitempty" yaml:"version,omitempty"`
}

// AWSBackend represents aws backend configuration
type AWSBackend struct {
	AccountID   *string `yaml:"account_id,omitempty"`
	AccountName string  `json:"account_name" yaml:"account_name"`
	Profile     string  `json:"profile" yaml:"profile"`
	Region      string  `json:"region" yaml:"region"`
	Bucket      string  `json:"bucket" yaml:"bucket"`
	DynamoTable *string `json:"dynamo_table" yaml:"dynamo_table"`
}

// Module is a module
type Module struct {
	Common `json:",inline" yaml:",inline"`
}

// Account is an account
type Account struct {
	ComponentCommon `json:",inline" yaml:",inline"`

	AllAccounts map[string]*json.Number `yaml:"all_accounts"`
	AccountName string                  `yaml:"account_name"`
	Global      *Component
}

// Component is a component
type Component struct {
	ComponentCommon `json:",inline" yaml:",inline"`

	Accounts  map[string]Account `json:"accounts" yaml:"accounts"` // Reference accounts for remote state
	Component string             `json:"component" yaml:"component"`
	EKS       *v1.EKSConfig      `json:"eks,omitempty" yaml:"eks,omitempty"`
	Env       string

	Kind            *v1.ComponentKind `json:"kind,omitempty" yaml:"kind,omitempty"`
	ModuleSource    *string           `json:"module_source" yaml:"module_source"`
	OtherComponents []string          `json:"other_components" yaml:"other_components"`
	Global          *Component        `json:"global" yaml:"global"`
}

// Env is an env
type Env struct {
	Components map[string]Component `json:"components" yaml:"components"`
	Env        string               `json:"env" yaml:"env"`
	EKS        *v1.EKSConfig        `json:"eks" yaml:"eks"`
}

// TfLint containts a plan for running tflint
type TfLint struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// Eval evaluates a config
func Eval(c *v2.Config) (*Plan, error) {
	if c == nil {
		return nil, errors.New("config is nil")
	}
	p := &Plan{}
	v, e := util.VersionString()
	if e != nil {
		return nil, errs.WrapInternal(e, "unable to parse fogg version")
	}
	p.Version = v

	var err error
	p.Global = p.buildGlobal(c)
	p.Accounts = p.buildAccounts(c)
	p.Envs, err = p.buildEnvs(c)
	if err != nil {
		return nil, err
	}

	p.Modules = p.buildModules(c)
	p.Atlantis = p.buildAtlantis()
	p.TravisCI = p.buildTravisCI(c, v)

	return p, nil
}

// Print prints a plan
func Print(p *Plan) error {
	out, err := yaml.Marshal(p)
	if err != nil {
		return errs.WrapInternal(err, "yaml: could not marshal")
	}
	fmt.Print(string(out))
	return nil
}

func (p *Plan) buildAccounts(c *v2.Config) map[string]Account {
	defaults := c.Defaults

	accountPlans := make(map[string]Account, len(c.Accounts))
	for name, acct := range c.Accounts {
		accountPlan := Account{}

		accountPlan.AccountName = name
		accountPlan.ComponentCommon = resolveComponentCommon(defaults.Common, acct.Common)
		accountPlan.AllAccounts = resolveAccounts(c.Accounts)
		accountPlan.PathToRepoRoot = "../../../"
		accountPlan.Global = &p.Global
		accountPlans[name] = accountPlan
	}

	return accountPlans
}

func (p *Plan) buildModules(c *v2.Config) map[string]Module {
	modulePlans := make(map[string]Module, len(c.Modules))
	for name, conf := range c.Modules {
		modulePlan := Module{}

		modulePlan.PathToRepoRoot = "../../../"
		modulePlan.TerraformVersion = *v2.ResolveModuleTerraformVersion(c.Defaults, conf)
		modulePlans[name] = modulePlan
	}
	return modulePlans
}

func newEnvPlan() Env {
	ep := Env{}
	ep.Components = make(map[string]Component)
	return ep
}

func (p *Plan) buildGlobal(conf *v2.Config) Component {
	// Global just uses defaults because that's the way sicc worked. We should make it directly configurable.
	componentPlan := Component{}
	componentPlan.Accounts = p.Accounts
	defaults := conf.Defaults
	global := conf.Global

	componentPlan.ComponentCommon = resolveComponentCommon(defaults.Common, global.Common)

	componentPlan.Component = "global"
	componentPlan.OtherComponents = []string{}
	componentPlan.ExtraVars = resolveExtraVars(defaults.ExtraVars, global.ExtraVars)
	componentPlan.PathToRepoRoot = "../../"

	return componentPlan
}

// buildEnvs must be build after accounts
func (p *Plan) buildEnvs(conf *v2.Config) (map[string]Env, error) {
	envPlans := make(map[string]Env, len(conf.Envs))
	defaults := conf.Defaults

	for envName, envConf := range conf.Envs {
		envPlan := newEnvPlan()
		envPlan.Env = envName

		for componentName, componentConf := range conf.Envs[envName].Components {
			componentPlan := Component{
				Kind: componentConf.Kind,
			}
			// reference accounts for remote state
			if _, dupe := p.Accounts[componentName]; dupe {
				return nil, errs.WrapUser(fmt.Errorf("Component %s can't have same name as account", componentName), "Invalid component name")
			}

			if componentConf.Kind.GetOrDefault() == v1.ComponentKindHelmTemplate {
				componentPlan.EKS = resolveEKSConfig(envPlan.EKS, componentConf.EKS)
			}
			componentPlan.Accounts = p.Accounts

			componentPlan.ComponentCommon = resolveComponentCommon(defaults.Common, envConf.Common, componentConf.Common)

			componentPlan.Env = envName
			componentPlan.Component = componentName
			componentPlan.OtherComponents = otherComponentNames(conf.Envs[envName].Components, componentName)
			componentPlan.ModuleSource = componentConf.ModuleSource
			componentPlan.PathToRepoRoot = "../../../../"

			componentPlan.Global = &p.Global

			envPlan.Components[componentName] = componentPlan
		}

		envPlans[envName] = envPlan
	}
	return envPlans, nil
}

func resolveComponentCommon(commons ...v2.Common) ComponentCommon {
	var awsPlan *AWSProvider
	awsConfig := v2.ResolveAWSProvider(commons...)

	if awsConfig != nil {
		awsPlan = &AWSProvider{
			AccountID:         *awsConfig.AccountID,
			Profile:           *awsConfig.Profile,
			Version:           *awsConfig.Version,
			Region:            *awsConfig.Region,
			AdditionalRegions: awsConfig.AdditionalRegions,
		}
	}

	var githubPlan *GithubProvider
	githubConfig := v2.ResolveGithubProvider(commons...)

	if githubConfig != nil {
		githubPlan = &GithubProvider {
			Organization: githubConfig.Organization,
			BaseURL: githubConfig.BaseURL,
		}
	}

	var snowflakePlan *SnowflakeProvider
	snowflakeConfig := v2.ResolveSnowflakeProvider(commons...)
	if snowflakeConfig != nil {
		snowflakePlan = &SnowflakeProvider{
			Account: *snowflakeConfig.Account,
			Role:    *snowflakeConfig.Role,
			Region:  *snowflakeConfig.Region,
			Version: snowflakeConfig.Version,
		}
	}

	var oktaPlan *OktaProvider
	oktaConfig := v2.ResolveOktaProvider(commons...)
	if oktaConfig != nil {
		oktaPlan = &OktaProvider{
			OrgName: *oktaConfig.OrgName,
			Version: oktaConfig.Version,
		}
	}

	var blessPlan *BlessProvider
	blessConfig := v2.ResolveBlessProvider(commons...)
	if blessConfig != nil && blessConfig.AWSProfile != nil && blessConfig.AWSRegion != nil {
		blessPlan = &BlessProvider{
			AWSProfile:        *blessConfig.AWSProfile,
			AWSRegion:         *blessConfig.AWSRegion,
			AdditionalRegions: blessConfig.AdditionalRegions,
			Version:           blessConfig.Version,
		}
	}

	tflintConfig := v2.ResolveTfLint(commons...)

	tfLintPlan := TfLint{
		Enabled: *tflintConfig.Enabled,
	}

	atlantisConfig := v2.ResolveAtlantis(commons...)

	atlantisPlan := AtlantisComponent{
		Enabled: *atlantisConfig.Enabled,
	}
	if atlantisPlan.Enabled {
		atlantisPlan.RoleName = *atlantisConfig.RoleName
		atlantisPlan.RolePath = *atlantisConfig.RolePath
	}

	travisConfig := v2.ResolveTravis(commons...)
	travisPlan := TravisComponent{
		Enabled: *travisConfig.Enabled,
	}
	if travisPlan.Enabled {
		travisPlan.AWSRoleName = *travisConfig.AWSIAMRoleName
		travisPlan.Command = *travisConfig.Command
	}

	return ComponentCommon{
		Atlantis: atlantisPlan,
		Backend: AWSBackend{
			AccountID:   v2.ResolveOptionalString(v2.BackendAccountIdGetter, commons...),
			Region:      v2.ResolveRequiredString(v2.BackendRegionGetter, commons...),
			Profile:     v2.ResolveRequiredString(v2.BackendProfileGetter, commons...),
			Bucket:      v2.ResolveRequiredString(v2.BackendBucketGetter, commons...),
			DynamoTable: v2.ResolveOptionalString(v2.BackendDynamoTableGetter, commons...),
		},
		Providers: Providers{
			AWS:       awsPlan,
			Github: githubPlan,
			Snowflake: snowflakePlan,
			Bless:     blessPlan,
			Okta:      oktaPlan,
		},
		TfLint:    tfLintPlan,
		ExtraVars: v2.ResolveStringMap(v2.ExtraVarsGetter, commons...),
		Owner:     v2.ResolveRequiredString(v2.OwnerGetter, commons...),
		Project:   v2.ResolveRequiredString(v2.ProjectGetter, commons...),
		Common:    Common{TerraformVersion: v2.ResolveRequiredString(v2.TerraformVersionGetter, commons...)},
		TravisCI:  travisPlan,
	}
}

func otherComponentNames(components map[string]v2.Component, thisComponent string) []string {
	r := make([]string, 0)
	for componentName, componentConf := range components {
		// Only set up remote state for terraform components
		if componentConf.Kind.GetOrDefault() != v1.ComponentKindTerraform {
			continue
		}
		if componentName != thisComponent {
			r = append(r, componentName)
		}
	}
	return r
}

func resolveEKSConfig(def *v1.EKSConfig, override *v1.EKSConfig) *v1.EKSConfig {
	resolved := &v1.EKSConfig{}
	if def != nil {
		resolved.ClusterName = def.ClusterName
	}
	if override != nil {
		resolved.ClusterName = override.ClusterName
	}
	return resolved
}

func resolveExtraVars(vars ...map[string]string) map[string]string {
	resolved := map[string]string{}

	for _, m := range vars {
		for k, v := range m {
			resolved[k] = v
		}
	}
	return resolved
}

func resolveAccounts(accounts map[string]v2.Account) map[string]*json.Number {
	a := make(map[string]*json.Number)
	for name, account := range accounts {
		var acctID *json.Number
		if account.Providers != nil && account.Providers.AWS != nil {
			acctID = account.Providers.AWS.AccountID
		}
		a[name] = acctID
	}
	return a
}
