package plan

import (
	"fmt"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/providers"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/pkg/errors"
)

type account struct {
	AccountID          *int64
	AccountName        string
	AllAccounts        map[string]int64
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSProviderVersion string
	AWSRegionBackend   string
	AWSRegionProvider  string
	AWSRegions         []string
	ExtraVars          map[string]string
	InfraBucket        string
	Owner              string
	Project            string
	TerraformVersion   string
}

type Module struct {
	TerraformVersion string
}

type Component struct {
	AccountID          *int64
	AccountName        string
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSProviderVersion string
	AWSRegionBackend   string
	AWSRegionProvider  string
	AWSRegions         []string
	Component          string
	Env                string
	ExtraVars          map[string]string
	InfraBucket        string
	ModuleSource       *string
	OtherComponents    []string
	Owner              string
	Project            string
	TerraformVersion   string
}

type Env struct {
	AccountID          *int64
	AccountName        string
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSProviderVersion string
	AWSRegionBackend   string
	AWSRegionProvider  string
	AWSRegions         []string
	Components         map[string]Component
	Env                string
	ExtraVars          map[string]string
	InfraBucket        string
	Owner              string
	Project            string
	TerraformVersion   string
	Type               string
}

type Plan struct {
	Accounts        map[string]account
	Envs            map[string]Env
	Global          Component
	Modules         map[string]Module
	Version         string
	CustomProviders map[string]*providers.CustomProvider
}

func Eval(config *config.Config, verbose bool) (*Plan, error) {
	p := &Plan{}
	v, e := util.VersionString()
	if e != nil {
		return nil, errors.Wrap(e, "unable to parse fogg version")
	}
	p.Version = v
	p.CustomProviders = config.Defaults.CustomProviders

	accounts, err := buildAccounts(config)
	if err != nil {
		return nil, err
	}
	p.Accounts = accounts

	envs, err := buildEnvs(config)
	if err != nil {
		return nil, err
	}
	p.Envs = envs

	global, err := buildGlobal(config)
	if err != nil {
		return nil, err
	}
	p.Global = global

	modules, err := buildModules(config)
	if err != nil {
		return nil, err
	}
	p.Modules = modules
	return p, nil
}

func Print(p *Plan) error {
	fmt.Printf("Version: %s\n", p.Version)
	fmt.Printf("fogg version: %s\n", p.Version)
	fmt.Println("Accounts:")
	for name, account := range p.Accounts {
		fmt.Printf("\t%s:\n", name)
		if account.AccountID != nil {
			fmt.Printf("\t\taccount id: %d\n", account.AccountID)
		}
		fmt.Printf("\t\taccount_id: %d\n", account.AccountID)

		fmt.Printf("\t\taws_profile_backend: %v\n", account.AWSProfileBackend)
		fmt.Printf("\t\taws_profile_provider: %v\n", account.AWSProfileProvider)
		fmt.Printf("\t\taws_provider_version: %v\n", account.AWSProviderVersion)
		fmt.Printf("\t\taws_region_backend: %v\n", account.AWSRegionBackend)
		fmt.Printf("\t\taws_region_provider: %v\n", account.AWSRegionProvider)
		fmt.Printf("\t\taws_regions: %v\n", account.AWSRegions)
		fmt.Printf("\t\tinfra_bucket: %v\n", account.InfraBucket)
		fmt.Printf("\t\tname: %v\n", account.AccountName)
		fmt.Printf("\t\towner: %v\n", account.Owner)
		fmt.Printf("\t\tproject: %v\n", account.Project)
		fmt.Printf("\t\tterraform_version: %v\n", account.TerraformVersion)

		fmt.Printf("\t\tall_accounts:\n")
		for acct, id := range account.AllAccounts {
			fmt.Printf("\t\t\t%s: %d\n", acct, id)
		}

	}

	fmt.Println("Global:")
	fmt.Printf("\taccount_id: %d\n", p.Global.AccountID)
	fmt.Printf("\taws_profile_backend: %v\n", p.Global.AWSProfileBackend)
	fmt.Printf("\taws_profile_provider: %v\n", p.Global.AWSProfileProvider)
	fmt.Printf("\taws_provider_version: %v\n", p.Global.AWSProviderVersion)
	fmt.Printf("\taws_region_backend: %v\n", p.Global.AWSRegionBackend)
	fmt.Printf("\taws_region_provider: %v\n", p.Global.AWSRegionProvider)
	fmt.Printf("\taws_regions: %v\n", p.Global.AWSRegions)
	fmt.Printf("\tinfra_bucket: %v\n", p.Global.InfraBucket)
	fmt.Printf("\tname: %v\n", p.Global.AccountName)
	fmt.Printf("\tother_p.Globals: %v\n", p.Global.OtherComponents)
	fmt.Printf("\towner: %v\n", p.Global.Owner)
	fmt.Printf("\tproject: %v\n", p.Global.Project)
	fmt.Printf("\tterraform_version: %v\n", p.Global.TerraformVersion)

	fmt.Println("Envs:")

	for name, env := range p.Envs {
		fmt.Printf("\t%s:\n", name)
		fmt.Printf("\t\taccount_id: %d\n", env.AccountID)

		fmt.Printf("\t\taws_profile_backend: %v\n", env.AWSProfileBackend)
		fmt.Printf("\t\taws_profile_provider: %v\n", env.AWSProfileProvider)
		fmt.Printf("\t\taws_provider_version: %v\n", env.AWSProviderVersion)
		fmt.Printf("\t\taws_region_backend: %v\n", env.AWSRegionBackend)
		fmt.Printf("\t\taws_region_provider: %v\n", env.AWSRegionProvider)
		fmt.Printf("\t\taws_regions: %v\n", env.AWSRegions)
		fmt.Printf("\t\tenv: %v\n", env.Env)
		fmt.Printf("\t\tinfra_bucket: %v\n", env.InfraBucket)
		fmt.Printf("\t\tname: %v\n", env.AccountName)
		fmt.Printf("\t\towner: %v\n", env.Owner)
		fmt.Printf("\t\tproject: %v\n", env.Project)
		fmt.Printf("\t\tterraform_version: %v\n", env.TerraformVersion)

		fmt.Println("\t\tComponents:")

		for name, component := range env.Components {
			fmt.Printf("\t\t\t%s:\n", name)
			fmt.Printf("\t\t\t\taccount_id: %d\n", component.AccountID)

			fmt.Printf("\t\t\t\taws_profile_backend: %v\n", component.AWSProfileBackend)
			fmt.Printf("\t\t\t\taws_profile_provider: %v\n", component.AWSProfileProvider)
			fmt.Printf("\t\t\t\taws_provider_version: %v\n", component.AWSProviderVersion)
			fmt.Printf("\t\t\t\taws_region_backend: %v\n", component.AWSRegionBackend)
			fmt.Printf("\t\t\t\taws_region_provider: %v\n", component.AWSRegionProvider)
			fmt.Printf("\t\t\t\taws_regions: %v\n", component.AWSRegions)
			fmt.Printf("\t\t\t\tinfra_bucket: %v\n", component.InfraBucket)
			fmt.Printf("\t\t\t\tname: %v\n", component.AccountName)
			fmt.Printf("\t\t\t\tother_components: %v\n", component.OtherComponents)
			fmt.Printf("\t\t\t\towner: %v\n", component.Owner)
			fmt.Printf("\t\t\t\tproject: %v\n", component.Project)
			fmt.Printf("\t\t\t\tterraform_version: %v\n", component.TerraformVersion)
		}

	}

	fmt.Println("Modules:")
	for name, module := range p.Modules {
		fmt.Printf("\t%s:\n", name)
		fmt.Printf("\t\tterraform_version: %s\n", module.TerraformVersion)
	}
	return nil
}

func buildAccounts(c *config.Config) (map[string]account, error) {
	defaults := c.Defaults

	accountPlans := make(map[string]account, len(c.Accounts))
	for name, config := range c.Accounts {
		accountPlan := account{}

		accountPlan.AccountName = name
		accountPlan.AccountID = resolveOptionalInt(c.Defaults.AccountID, config.AccountID)

		accountPlan.AWSRegionBackend = resolveRequired(defaults.AWSRegionBackend, config.AWSRegionBackend)
		accountPlan.AWSRegionProvider = resolveRequired(defaults.AWSRegionProvider, config.AWSRegionProvider)
		accountPlan.AWSRegions = resolveStringArray(defaults.AWSRegions, config.AWSRegions)

		accountPlan.AWSProfileBackend = resolveRequired(defaults.AWSProfileBackend, config.AWSProfileBackend)
		accountPlan.AWSProfileProvider = resolveRequired(defaults.AWSProfileProvider, config.AWSProfileProvider)
		accountPlan.AWSProviderVersion = resolveRequired(defaults.AWSProviderVersion, config.AWSProviderVersion)
		accountPlan.AllAccounts = resolveAccounts(c.Accounts)
		accountPlan.TerraformVersion = resolveRequired(defaults.TerraformVersion, config.TerraformVersion)
		accountPlan.InfraBucket = resolveRequired(defaults.InfraBucket, config.InfraBucket)
		accountPlan.Owner = resolveRequired(defaults.Owner, config.Owner)
		accountPlan.Project = resolveRequired(defaults.Project, config.Project)
		accountPlan.ExtraVars = resolveExtraVars(defaults.ExtraVars, config.ExtraVars)

		accountPlans[name] = accountPlan
	}

	return accountPlans, nil
}

func buildModules(c *config.Config) (map[string]Module, error) {
	modulePlans := make(map[string]Module, len(c.Modules))
	for name, conf := range c.Modules {
		modulePlan := Module{}

		modulePlan.TerraformVersion = resolveRequired(c.Defaults.TerraformVersion, conf.TerraformVersion)
		modulePlans[name] = modulePlan
	}
	return modulePlans, nil
}

func newEnvPlan() Env {
	ep := Env{}
	ep.Components = make(map[string]Component)
	return ep
}

func buildGlobal(conf *config.Config) (Component, error) {
	// Global just uses defaults because that's the way sicc worked. We should make it directly configurable.
	componentPlan := Component{}

	componentPlan.AccountID = conf.Defaults.AccountID

	componentPlan.AWSRegionBackend = conf.Defaults.AWSRegionBackend
	componentPlan.AWSRegionProvider = conf.Defaults.AWSRegionProvider
	componentPlan.AWSRegions = conf.Defaults.AWSRegions

	componentPlan.AWSProfileBackend = conf.Defaults.AWSProfileBackend
	componentPlan.AWSProfileProvider = conf.Defaults.AWSProfileProvider
	componentPlan.AWSProviderVersion = conf.Defaults.AWSProviderVersion
	// TODO add AccountID to defaults
	// componentPlan.AccountID = conf.Defaults.AccountID

	componentPlan.TerraformVersion = conf.Defaults.TerraformVersion
	componentPlan.InfraBucket = conf.Defaults.InfraBucket
	componentPlan.Owner = conf.Defaults.Owner
	componentPlan.Project = conf.Defaults.Project
	componentPlan.ExtraVars = conf.Defaults.ExtraVars

	componentPlan.Component = "global"
	return componentPlan, nil
}

func buildEnvs(conf *config.Config) (map[string]Env, error) {
	envPlans := make(map[string]Env, len(conf.Envs))
	defaults := conf.Defaults

	defaultExtraVars := defaults.ExtraVars

	for envName, envConf := range conf.Envs {
		envPlan := newEnvPlan()

		envPlan.AccountID = resolveOptionalInt(conf.Defaults.AccountID, envConf.AccountID)
		envPlan.Env = envName

		envPlan.AWSRegionBackend = resolveRequired(defaults.AWSRegionBackend, envConf.AWSRegionBackend)
		envPlan.AWSRegionProvider = resolveRequired(defaults.AWSRegionProvider, envConf.AWSRegionProvider)
		envPlan.AWSRegions = resolveStringArray(defaults.AWSRegions, envConf.AWSRegions)

		envPlan.AWSProfileBackend = resolveRequired(defaults.AWSProfileBackend, envConf.AWSProfileBackend)
		envPlan.AWSProfileProvider = resolveRequired(defaults.AWSProfileProvider, envConf.AWSProfileProvider)
		envPlan.AWSProviderVersion = resolveRequired(defaults.AWSProviderVersion, envConf.AWSProviderVersion)

		envPlan.TerraformVersion = resolveRequired(defaults.TerraformVersion, envConf.TerraformVersion)
		envPlan.InfraBucket = resolveRequired(defaults.InfraBucket, envConf.InfraBucket)
		envPlan.Owner = resolveRequired(defaults.Owner, envConf.Owner)
		envPlan.Project = resolveRequired(defaults.Project, envConf.Project)
		envPlan.ExtraVars = resolveExtraVars(defaultExtraVars, envConf.ExtraVars)

		for componentName, componentConf := range conf.Envs[envName].Components {
			componentPlan := Component{}

			componentPlan.AccountID = resolveOptionalInt(envPlan.AccountID, componentConf.AccountID)
			componentPlan.AWSRegionBackend = resolveRequired(envPlan.AWSRegionBackend, componentConf.AWSRegionBackend)
			componentPlan.AWSRegionProvider = resolveRequired(envPlan.AWSRegionProvider, componentConf.AWSRegionProvider)
			componentPlan.AWSRegions = resolveStringArray(envPlan.AWSRegions, componentConf.AWSRegions)

			componentPlan.AWSProfileBackend = resolveRequired(envPlan.AWSProfileBackend, componentConf.AWSProfileBackend)
			componentPlan.AWSProfileProvider = resolveRequired(envPlan.AWSProfileProvider, componentConf.AWSProfileProvider)
			componentPlan.AWSProviderVersion = resolveRequired(envPlan.AWSProviderVersion, componentConf.AWSProviderVersion)
			componentPlan.AccountID = resolveOptionalInt(envPlan.AccountID, componentConf.AccountID)

			componentPlan.TerraformVersion = resolveRequired(envPlan.TerraformVersion, componentConf.TerraformVersion)
			componentPlan.InfraBucket = resolveRequired(envPlan.InfraBucket, componentConf.InfraBucket)
			componentPlan.Owner = resolveRequired(envPlan.Owner, componentConf.Owner)
			componentPlan.Project = resolveRequired(envPlan.Project, componentConf.Project)

			componentPlan.Env = envName
			componentPlan.Component = componentName
			componentPlan.OtherComponents = otherComponentNames(conf.Envs[envName].Components, componentName)
			// This is a bit awkward but should go away when we make the modules thing first-class.
			if envPlan.Type == "aws" {
				componentPlan.OtherComponents = append(componentPlan.OtherComponents, "cloud-env")
			}
			componentPlan.ModuleSource = componentConf.ModuleSource
			componentPlan.ExtraVars = resolveExtraVars(envPlan.ExtraVars, componentConf.ExtraVars)

			envPlan.Components[componentName] = componentPlan
		}

		envPlans[envName] = envPlan
	}
	return envPlans, nil
}

func otherComponentNames(components map[string]*config.Component, thisComponent string) []string {
	r := make([]string, 0)
	for componentName := range components {
		if componentName != thisComponent {
			r = append(r, componentName)
		}
	}
	return r
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

func resolveOptionalInt(def *int64, override *int64) *int64 {
	if override != nil {
		return override
	}
	return def
}

func resolveAccounts(accounts map[string]config.Account) map[string]int64 {
	a := make(map[string]int64)
	for name, account := range accounts {
		if account.AccountID != nil {
			a[name] = *account.AccountID
		}
	}
	return a
}
