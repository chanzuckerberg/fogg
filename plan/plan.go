package plan

import (
	"fmt"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
)

type account struct {
	AccountId          *int64
	AccountName        string
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSRegion          string
	AWSRegions         []string
	InfraBucket        string
	OtherAccounts      map[string]int64
	Owner              string
	Project            string
	TerraformVersion   string
}

type module struct {
	TerraformVersion string
}

type component struct {
	AccountId          *int64
	AccountName        string
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSRegion          string
	AWSRegions         []string
	InfraBucket        string
	Owner              string
	Project            string
	Env                string
	TerraformVersion   string
}

type env struct {
	AccountId          *int64
	AccountName        string
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSRegion          string
	AWSRegions         []string
	InfraBucket        string
	Owner              string
	Project            string
	Env                string
	TerraformVersion   string

	Components map[string]*component
}

type plan struct {
	Accounts map[string]*account
	Version  string
	Modules  map[string]*module
	Envs     map[string]*env
}

func Plan(fs afero.Fs, configFile string) (*plan, error) {
	c, err := config.FindAndReadConfig(fs, configFile)
	if err != nil {
		return nil, err
	}
	p := &plan{}
	// read config and validate
	// build repo plan
	p.Version = util.VersionString()
	p.Accounts = buildAccounts(c)
	p.Envs = buildEnvs(c)
	p.Modules = buildModules(c)
	return p, nil
}

func Print(p *plan) error {
	fmt.Printf("Version: %s\n", p.Version)
	fmt.Println("Accounts:")
	for name, account := range p.Accounts {
		fmt.Printf("\t%s:\n", name)
		if account.AccountId != nil {
			fmt.Printf("\t\taccount id: %d\n", account.AccountId)
		}
		fmt.Printf("\t\tid: %d\n", account.AccountId)
		fmt.Printf("\t\tname: %v\n", account.AccountName)
		fmt.Printf("\t\taws_profile_backend: %v\n", account.AWSProfileBackend)
		fmt.Printf("\t\taws_profile_provider: %v\n", account.AWSProfileProvider)
		fmt.Printf("\t\taws_region: %v\n", account.AWSRegion)
		fmt.Printf("\t\taws_regions: %v\n", account.AWSRegions)
		fmt.Printf("\t\tinfra_bucket: %v\n", account.InfraBucket)
		fmt.Printf("\t\towner: %v\n", account.Owner)
		fmt.Printf("\t\tproject: %v\n", account.Project)
		fmt.Printf("\t\tterraform_version: %v\n", account.TerraformVersion)
		fmt.Printf("\t\tother_accounts:\n")
		for acct, id := range account.OtherAccounts {
			fmt.Printf("\t\t\t%s: %d\n", acct, id)
		}

	}

	fmt.Println("Envs:")

	for name, env := range p.Envs {
		fmt.Printf("%s:\n", name)
	}

	fmt.Println("Modules:")
	for name, module := range p.Modules {
		fmt.Printf("%s:\n", name)
		fmt.Printf("\tterraform_version: %s\n", module.TerraformVersion)
	}
	return nil
}

func buildAccounts(c *config.Config) map[string]*account {
	defaults := c.Defaults
	accountPlans := make(map[string]*account, len(c.Accounts))
	for name, config := range c.Accounts {
		accountPlan := &account{}

		accountPlan.AccountName = name
		accountPlan.AccountId = config.AccountId

		accountPlan.AWSRegion = resolveRequired(defaults.AWSRegion, config.AWSRegion)
		accountPlan.AWSRegions = resolveStringArray(defaults.AWSRegions, config.AWSRegions)

		accountPlan.AWSProfileBackend = resolveRequired(defaults.AWSProfileBackend, config.AWSProfileBackend)
		accountPlan.AWSProfileProvider = resolveRequired(defaults.AWSProfileProvider, config.AWSProfileProvider)
		accountPlan.OtherAccounts = resolveOtherAccounts(c.Accounts, name)
		accountPlan.TerraformVersion = resolveRequired(defaults.TerraformVersion, config.TerraformVersion)
		accountPlan.InfraBucket = resolveRequired(defaults.InfraBucket, config.InfraBucket)
		accountPlan.Owner = resolveRequired(defaults.Owner, config.Owner)
		accountPlan.Project = resolveRequired(defaults.Project, config.Project)

		accountPlans[name] = accountPlan
	}

	return accountPlans
}

func buildModules(c *config.Config) map[string]*module {
	modulePlans := make(map[string]*module, len(c.Modules))
	for name, conf := range c.Modules {
		modulePlan := &module{}

		modulePlan.TerraformVersion = resolveRequired(c.Defaults.TerraformVersion, conf.TerraformVersion)
		modulePlans[name] = modulePlan
	}
	return modulePlans
}

func newEnvPlan() *env {
	ep := &env{}
	ep.Components = make(map[string]*component)
	return ep
}

func buildEnvs(conf *config.Config) map[string]*env {
	envPlans := make(map[string]*env, len(conf.Envs))
	defaults := conf.Defaults
	for envName, envConf := range conf.Envs {
		envPlan := newEnvPlan()

		envPlan.AccountId = envConf.AccountId

		envPlan.AWSRegion = resolveRequired(defaults.AWSRegion, envConf.AWSRegion)
		envPlan.AWSRegions = resolveStringArray(defaults.AWSRegions, envConf.AWSRegions)

		envPlan.AWSProfileBackend = resolveRequired(defaults.AWSProfileBackend, envConf.AWSProfileBackend)
		envPlan.AWSProfileProvider = resolveRequired(defaults.AWSProfileProvider, envConf.AWSProfileProvider)

		envPlan.TerraformVersion = resolveRequired(defaults.TerraformVersion, envConf.TerraformVersion)
		envPlan.InfraBucket = resolveRequired(defaults.InfraBucket, envConf.InfraBucket)
		envPlan.Owner = resolveRequired(defaults.Owner, envConf.Owner)
		envPlan.Project = resolveRequired(defaults.Project, envConf.Project)

		components := make(map[string]*component)

		// FIXME no longer needed
		for name, _ := range conf.Envs[envName].Components {
			components[name] = nil
		}

		for componentName, _ := range components {
			componentPlan := &component{}
			componentConf := envConf
			componentPlan.AccountId = resolveOptionalInt(envPlan.AccountId, componentConf.AccountId)

			componentPlan.AWSRegion = resolveRequired(envPlan.AWSRegion, componentConf.AWSRegion)
			componentPlan.AWSRegions = resolveStringArray(&envPlan.AWSRegions, componentConf.AWSRegions)

			componentPlan.AWSProfileBackend = resolveRequired(envPlan.AWSProfileBackend, componentConf.AWSProfileBackend)
			componentPlan.AWSProfileProvider = resolveRequired(envPlan.AWSProfileProvider, componentConf.AWSProfileProvider)

			componentPlan.TerraformVersion = resolveRequired(envPlan.TerraformVersion, componentConf.TerraformVersion)
			componentPlan.InfraBucket = resolveRequired(envPlan.InfraBucket, componentConf.InfraBucket)
			componentPlan.Owner = resolveRequired(envPlan.Owner, componentConf.Owner)
			componentPlan.Project = resolveRequired(envPlan.Project, componentConf.Project)

			envPlan.Components[componentName] = componentPlan
		}

		envPlans[envName] = envPlan
	}
	return envPlans
}

func resolveStringArray(def *[]string, override *[]string) []string {
	if override != nil {
		return *override
	}
	if def != nil {
		return *def
	}
	return []string{}
}

func resolveRequired(def string, override *string) string {
	if override != nil {
		return *override
	}
	return def
}

func resolveOptional(def *string, override *string) *string {
	if override != nil {
		return override
	}
	return def
}

func resolveOptionalInt(def *int64, override *int64) *int64 {
	if override != nil {
		return override
	}
	return def
}

func resolveOtherAccounts(accounts map[string]config.Account, currentAccount string) map[string]int64 {
	other := make(map[string]int64)
	for name, account := range accounts {
		if name != currentAccount && account.AccountId != nil {
			other[name] = *account.AccountId
		}
	}
	return other
}
