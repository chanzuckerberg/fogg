package plan

import (
	"fmt"

	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/util"
	yaml "gopkg.in/yaml.v2"
)

const (
	// The version of the chanzuckerberg/terraform docker image to use
	dockerImageVersion = "0.2.1"
)

// AWSConfiguration represents aws configuration
type AWSConfiguration struct {
	AccountID          int64    `yaml:"account_id"`
	AccountName        string   `yaml:"account_name"`
	AWSProfileBackend  string   `yaml:"aws_profile_backend"`
	AWSProfileProvider string   `yaml:"aws_profile_provider"`
	AWSProviderVersion string   `yaml:"aws_provider_version"`
	AWSRegionBackend   string   `yaml:"aws_region_backend"`
	AWSRegionProvider  string   `yaml:"aws_region_provider"`
	AWSRegions         []string `yaml:"aws_regions"`
	InfraBucket        string   `yaml:"infra_bucket"`
	InfraDynamoTable   string   `yaml:"infra_dynamo_table"`
}

// Common represents common fields
type Common struct {
	Docker             bool   `yaml:"docker"`
	DockerImageVersion string `yaml:"docker_image_version"`
	PathToRepoRoot     string `yaml:"path_to_repo_root"`
	TerraformVersion   string `yaml:"terraform_version"`
}

// Account is an account
type Account struct {
	AWSConfiguration `yaml:",inline"`
	Common           `yaml:",inline"`

	AllAccounts map[string]int64  `yaml:"all_accounts"`
	ExtraVars   map[string]string `yaml:"extra_vars"`
	Owner       string
	Project     string
	TfLint      TfLint
}

// Module is a module
type Module struct {
	Common `yaml:",inline"`
}

// Component is a component
type Component struct {
	Accounts         map[string]Account // Reference accounts for remote state
	AWSConfiguration `yaml:",inline"`
	Common           `yaml:",inline"`
	Component        string
	EKS              *v1.EKSConfig `yaml:"eks,omitempty"`
	Env              string
	ExtraVars        map[string]string `yaml:"extra_vars"`
	Kind             *v1.ComponentKind `yaml:"kind,omitempty"`
	ModuleSource     *string           `yaml:"module_source"`
	OtherComponents  []string          `yaml:"other_components"`
	Owner            string
	Project          string
	TfLint           TfLint
}

// Env is an env
type Env struct {
	AWSConfiguration `yaml:",inline"`
	Common           `yaml:",inline"`

	Components map[string]Component
	Env        string
	EKS        *v1.EKSConfig
	ExtraVars  map[string]string `yaml:"extra_vars"`
	Owner      string
	Project    string
	TfLint     TfLint
}

// TfLint containts a plan for running tflint
type TfLint struct {
	Enabled bool
}

type AWSProfile struct {
	Name string
	ID   int64
	Role string
}

// Plan represents a set of actions to take
type Plan struct {
	Accounts map[string]Account
	Envs     map[string]Env
	Global   Component
	Modules  map[string]Module
	TravisCI TravisCI
	Version  string
}

// Eval evaluates a config
func Eval(config *v1.Config) (*Plan, error) {
	p := &Plan{}
	v, e := util.VersionString()
	if e != nil {
		return nil, errs.WrapInternal(e, "unable to parse fogg version")
	}
	p.Version = v

	var err error
	p.Accounts = p.buildAccounts(config)
	p.Envs, err = p.buildEnvs(config)
	if err != nil {
		return nil, err
	}
	p.Global = p.buildGlobal(config)
	p.Modules = p.buildModules(config)

	if config.TravisCI != nil {
		p.TravisCI = p.buildTravisCI(config, v)
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

func (p *Plan) buildAccounts(c *v1.Config) map[string]Account {
	defaults := c.Defaults

	accountPlans := make(map[string]Account, len(c.Accounts))
	for name, config := range c.Accounts {
		accountPlan := Account{}
		accountPlan.DockerImageVersion = dockerImageVersion

		accountPlan.AccountName = name
		accountPlan.AccountID = resolveRequiredInt(c.Defaults.AccountID, config.AccountID)

		accountPlan.AWSRegionBackend = resolveRequired(defaults.AWSRegionBackend, config.AWSRegionBackend)
		accountPlan.AWSRegionProvider = resolveRequired(defaults.AWSRegionProvider, config.AWSRegionProvider)
		accountPlan.AWSRegions = resolveStringArray(defaults.AWSRegions, config.AWSRegions)

		accountPlan.AWSProfileBackend = resolveRequired(defaults.AWSProfileBackend, config.AWSProfileBackend)
		accountPlan.AWSProfileProvider = resolveRequired(defaults.AWSProfileProvider, config.AWSProfileProvider)
		accountPlan.AWSProviderVersion = resolveRequired(defaults.AWSProviderVersion, config.AWSProviderVersion)
		accountPlan.AllAccounts = resolveAccounts(c.Accounts)
		accountPlan.TerraformVersion = resolveRequired(defaults.TerraformVersion, config.TerraformVersion)
		accountPlan.InfraBucket = resolveRequired(defaults.InfraBucket, config.InfraBucket)
		accountPlan.InfraDynamoTable = resolveRequired(defaults.InfraDynamoTable, config.InfraDynamoTable)
		accountPlan.Owner = resolveRequired(defaults.Owner, config.Owner)
		accountPlan.PathToRepoRoot = "../../../"
		accountPlan.Project = resolveRequired(defaults.Project, config.Project)
		accountPlan.ExtraVars = resolveExtraVars(defaults.ExtraVars, config.ExtraVars)

		accountPlan.TfLint = resolveTfLint(defaults.TfLint, config.TfLint)
		accountPlan.Docker = c.Docker

		accountPlans[name] = accountPlan
	}

	return accountPlans
}

func (p *Plan) buildModules(c *v1.Config) map[string]Module {
	modulePlans := make(map[string]Module, len(c.Modules))
	for name, conf := range c.Modules {
		modulePlan := Module{}

		modulePlan.DockerImageVersion = dockerImageVersion
		modulePlan.PathToRepoRoot = "../../../"
		modulePlan.TerraformVersion = resolveRequired(c.Defaults.TerraformVersion, conf.TerraformVersion)
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

func (p *Plan) buildGlobal(conf *v1.Config) Component {
	// Global just uses defaults because that's the way sicc worked. We should make it directly configurable.
	componentPlan := Component{}

	componentPlan.DockerImageVersion = dockerImageVersion
	componentPlan.AccountID = conf.Defaults.AccountID

	componentPlan.AWSRegionBackend = conf.Defaults.AWSRegionBackend
	componentPlan.AWSRegionProvider = conf.Defaults.AWSRegionProvider
	componentPlan.AWSRegions = conf.Defaults.AWSRegions

	componentPlan.AWSProfileBackend = conf.Defaults.AWSProfileBackend
	componentPlan.AWSProfileProvider = conf.Defaults.AWSProfileProvider
	componentPlan.AWSProviderVersion = conf.Defaults.AWSProviderVersion

	componentPlan.TerraformVersion = conf.Defaults.TerraformVersion
	componentPlan.InfraBucket = conf.Defaults.InfraBucket
	componentPlan.InfraDynamoTable = conf.Defaults.InfraDynamoTable
	componentPlan.Owner = conf.Defaults.Owner
	componentPlan.PathToRepoRoot = "../../"
	componentPlan.Project = conf.Defaults.Project
	componentPlan.ExtraVars = conf.Defaults.ExtraVars
	componentPlan.TfLint = resolveTfLint(conf.Defaults.TfLint, nil)
	componentPlan.Docker = conf.Docker
	componentPlan.Component = "global"
	return componentPlan
}

// buildEnvs must be build after accounts
func (p *Plan) buildEnvs(conf *v1.Config) (map[string]Env, error) {
	envPlans := make(map[string]Env, len(conf.Envs))
	defaults := conf.Defaults

	defaultExtraVars := defaults.ExtraVars

	for envName, envConf := range conf.Envs {
		envPlan := newEnvPlan()

		envPlan.AccountID = resolveRequiredInt(conf.Defaults.AccountID, envConf.AccountID)
		envPlan.Env = envName
		envPlan.DockerImageVersion = dockerImageVersion

		envPlan.AWSRegionBackend = resolveRequired(defaults.AWSRegionBackend, envConf.AWSRegionBackend)
		envPlan.AWSRegionProvider = resolveRequired(defaults.AWSRegionProvider, envConf.AWSRegionProvider)
		envPlan.AWSRegions = resolveStringArray(defaults.AWSRegions, envConf.AWSRegions)

		envPlan.AWSProfileBackend = resolveRequired(defaults.AWSProfileBackend, envConf.AWSProfileBackend)
		envPlan.AWSProfileProvider = resolveRequired(defaults.AWSProfileProvider, envConf.AWSProfileProvider)
		envPlan.AWSProviderVersion = resolveRequired(defaults.AWSProviderVersion, envConf.AWSProviderVersion)

		envPlan.TerraformVersion = resolveRequired(defaults.TerraformVersion, envConf.TerraformVersion)
		envPlan.InfraBucket = resolveRequired(defaults.InfraBucket, envConf.InfraBucket)
		envPlan.InfraDynamoTable = resolveRequired(defaults.InfraDynamoTable, envConf.InfraDynamoTable)
		envPlan.Owner = resolveRequired(defaults.Owner, envConf.Owner)
		envPlan.Project = resolveRequired(defaults.Project, envConf.Project)
		envPlan.ExtraVars = resolveExtraVars(defaultExtraVars, envConf.ExtraVars)
		envPlan.TfLint = resolveTfLint(defaults.TfLint, envConf.TfLint)
		envPlan.Docker = conf.Docker

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

			componentPlan.AccountID = resolveRequiredInt(envPlan.AccountID, componentConf.AccountID)
			componentPlan.AWSRegionBackend = resolveRequired(envPlan.AWSRegionBackend, componentConf.AWSRegionBackend)
			componentPlan.AWSRegionProvider = resolveRequired(envPlan.AWSRegionProvider, componentConf.AWSRegionProvider)
			componentPlan.AWSRegions = resolveStringArray(envPlan.AWSRegions, componentConf.AWSRegions)

			componentPlan.AWSProfileBackend = resolveRequired(envPlan.AWSProfileBackend, componentConf.AWSProfileBackend)
			componentPlan.AWSProfileProvider = resolveRequired(envPlan.AWSProfileProvider, componentConf.AWSProfileProvider)
			componentPlan.AWSProviderVersion = resolveRequired(envPlan.AWSProviderVersion, componentConf.AWSProviderVersion)

			componentPlan.TerraformVersion = resolveRequired(envPlan.TerraformVersion, componentConf.TerraformVersion)
			componentPlan.InfraBucket = resolveRequired(envPlan.InfraBucket, componentConf.InfraBucket)
			componentPlan.InfraDynamoTable = resolveRequired(envPlan.InfraDynamoTable, componentConf.InfraDynamoTable)
			componentPlan.Owner = resolveRequired(envPlan.Owner, componentConf.Owner)
			componentPlan.Project = resolveRequired(envPlan.Project, componentConf.Project)

			componentPlan.Env = envName
			componentPlan.Component = componentName
			componentPlan.DockerImageVersion = dockerImageVersion
			componentPlan.OtherComponents = otherComponentNames(conf.Envs[envName].Components, componentName)
			componentPlan.ModuleSource = componentConf.ModuleSource
			componentPlan.ExtraVars = resolveExtraVars(envPlan.ExtraVars, componentConf.ExtraVars)
			componentPlan.PathToRepoRoot = "../../../../"

			componentPlan.TfLint = resolveTfLintComponent(envPlan.TfLint, componentConf.TfLint)
			componentPlan.Docker = conf.Docker

			envPlan.Components[componentName] = componentPlan
		}

		envPlans[envName] = envPlan
	}
	return envPlans, nil
}

func otherComponentNames(components map[string]*v1.Component, thisComponent string) []string {
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

func resolveExtraVars(def map[string]string, override map[string]string) map[string]string {
	resolved := map[string]string{}
	for k, v := range def {
		resolved[k] = v
	}
	for k, v := range override {
		resolved[k] = v
	}
	return resolved
}

func resolveStringArray(def []string, override []string) []string {
	if override != nil {
		return override
	}
	return def
}

func resolveRequired(def string, override *string) string {
	if override != nil {
		return *override
	}
	return def
}

func resolveRequiredInt(def int64, override *int64) int64 {
	if override != nil {
		return *override
	}
	return def
}

func resolveOptionalInt(def *int64, override *int64) *int64 {
	if override != nil {
		return override
	}
	return def
}

func resolveAccounts(accounts map[string]v1.Account) map[string]int64 {
	a := make(map[string]int64)
	for name, account := range accounts {
		if account.AccountID != nil {
			a[name] = *account.AccountID
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
