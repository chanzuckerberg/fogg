package plan

import (
	"errors"
	"fmt"

	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/util"
	yaml "gopkg.in/yaml.v2"
)

const (
	// The version of the chanzuckerberg/terraform docker image to use
	dockerImageVersion = "0.2.1"
)

// Plan represents a set of actions to take
type Plan struct {
	Accounts map[string]Account
	Envs     map[string]Env
	Global   Component
	Modules  map[string]Module
	TravisCI TravisCI
	Version  string
}

// Common represents common fields
type Common struct {
	Docker             bool   `yaml:"docker"`
	DockerImageVersion string `yaml:"docker_image_version"`
	PathToRepoRoot     string `yaml:"path_to_repo_root"`
	TerraformVersion   string `yaml:"terraform_version"`
}

type ComponentCommon struct {
	Common `yaml:",inline"`

	Backend   AWSBackend        `yaml:"backend"`
	ExtraVars map[string]string `yaml:"extra_vars"`
	Owner     string
	Project   string
	Providers Providers `yaml:"providers"`
	TfLint    TfLint
}

type Providers struct {
	AWS *AWSProvider `yaml:"aws"`
}

type AWSProvider struct {
	AccountID         int64    `yaml:"account_id"`
	Profile           string   `yaml:"profile"`
	Version           string   `yaml:"version"`
	Region            string   `yaml:"region"`
	AdditionalRegions []string `yaml:"additional_regions"`
}

// AWSConfiguration represents aws configuration
type AWSBackend struct {
	AccountName string  `yaml:"account_name"`
	Profile     string  `yaml:"profile"`
	Region      string  `yaml:"region"`
	Bucket      string  `yaml:"bucket"`
	DynamoTable *string `yaml:"dynamo_table"`
}

// Module is a module
type Module struct {
	Common `yaml:",inline"`
}

// Account is an account
type Account struct {
	ComponentCommon `yaml:",inline"`

	AllAccounts map[string]int64 `yaml:"all_accounts"`
	AccountName string           `yaml:"account_name"`
}

// Component is a component
type Component struct {
	ComponentCommon `yaml:",inline"`

	Accounts  map[string]Account // Reference accounts for remote state
	Component string
	EKS       *v1.EKSConfig `yaml:"eks,omitempty"`
	Env       string

	Kind            *v1.ComponentKind `yaml:"kind,omitempty"`
	ModuleSource    *string           `yaml:"module_source"`
	OtherComponents []string          `yaml:"other_components"`
}

// Env is an env
type Env struct {
	Components map[string]Component
	Env        string
	EKS        *v1.EKSConfig //TODO get rid of this
}

// TfLint containts a plan for running tflint
type TfLint struct {
	Enabled bool
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
	p.Accounts = p.buildAccounts(c)
	p.Envs, err = p.buildEnvs(c)
	if err != nil {
		return nil, err
	}
	p.Global = p.buildGlobal(c)
	p.Modules = p.buildModules(c)

	if c.Tools.TravisCI != nil {
		p.TravisCI = p.buildTravisCI(c, v)
	}

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
		accountPlan.TfLint = resolveTfLint(c.Tools.TfLint, nil)
		accountPlan.AllAccounts = resolveAccounts(c.Accounts)
		accountPlan.PathToRepoRoot = "../../../"
		accountPlan.Docker = c.Docker
		accountPlan.DockerImageVersion = dockerImageVersion

		accountPlans[name] = accountPlan
	}

	return accountPlans
}

func (p *Plan) buildModules(c *v2.Config) map[string]Module {
	modulePlans := make(map[string]Module, len(c.Modules))
	for name, conf := range c.Modules {
		modulePlan := Module{}

		modulePlan.DockerImageVersion = dockerImageVersion
		modulePlan.PathToRepoRoot = "../../../"
		modulePlan.TerraformVersion = *v2.ResolveModuleTerraformVersion(c.Defaults, conf)
		modulePlan.Docker = c.Docker
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
	componentPlan.DockerImageVersion = dockerImageVersion
	componentPlan.OtherComponents = []string{}
	componentPlan.ExtraVars = resolveExtraVars(defaults.ExtraVars, global.ExtraVars)
	componentPlan.PathToRepoRoot = "../../"

	componentPlan.TfLint = resolveTfLint(conf.Tools.TfLint, nil)

	componentPlan.Docker = conf.Docker

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
			componentPlan.DockerImageVersion = dockerImageVersion
			componentPlan.OtherComponents = otherComponentNames(conf.Envs[envName].Components, componentName)
			componentPlan.ModuleSource = componentConf.ModuleSource
			componentPlan.PathToRepoRoot = "../../../../"

			componentPlan.TfLint = resolveTfLint(conf.Tools.TfLint, nil)
			componentPlan.Docker = conf.Docker

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

	return ComponentCommon{
		Backend: AWSBackend{
			Region:      v2.ResolveRequiredString(v2.BackendRegionGetter, commons...),
			Profile:     v2.ResolveRequiredString(v2.BackendProfileGetter, commons...),
			Bucket:      v2.ResolveRequiredString(v2.BackendBucketGetter, commons...),
			DynamoTable: v2.ResolveOptionalString(v2.BackendDynamoTableGetter, commons...),
		},
		Providers: Providers{
			AWS: awsPlan,
		},
		ExtraVars: v2.ResolveStringMap(v2.ExtraVarsGetter, commons...),
		Owner:     v2.ResolveRequiredString(v2.OwnerGetter, commons...),
		Project:   v2.ResolveRequiredString(v2.ProjectGetter, commons...),
		Common:    Common{TerraformVersion: v2.ResolveRequiredString(v2.TerraformVersionGetter, commons...)},
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

func resolveAccounts(accounts map[string]v2.Account) map[string]int64 {
	a := make(map[string]int64)
	for name, account := range accounts {
		if account.Providers != nil && account.Providers.AWS != nil && account.Providers.AWS.AccountID != nil {
			a[name] = *account.Providers.AWS.AccountID
		}
	}
	return a
}

func resolveTfLint(def *v1.TfLint, override *v1.TfLint) TfLint {
	enabled := false
	if def != nil && def.Enabled != nil {
		enabled = *def.Enabled
	}
	if override != nil && override.Enabled != nil {
		enabled = *override.Enabled
	}
	return TfLint{
		Enabled: enabled,
	}
}

func resolveTfLintComponent(def TfLint, override *v1.TfLint) TfLint {

	enabled := def.Enabled
	if override != nil && override.Enabled != nil {
		enabled = *override.Enabled
	}
	return TfLint{
		Enabled: enabled,
	}
}
