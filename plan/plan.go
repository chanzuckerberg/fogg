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

type plan struct {
	Accounts map[string]*account
	Version  string
	Modules  map[string]*module
}

func Plan(fs afero.Fs, configFile string) (*plan, error) {
	c, _ := config.FindAndReadConfig(fs, configFile)
	p := &plan{}
	// read config and validate
	// build repo plan
	// build .sicc version plan
	p.Version = util.VersionString()
	p.Accounts = buildAccounts(c)
	p.Modules = buildModules(c)
	// build modules plan

	// build envs plan
	// walk config and apply inheritance rules
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

		profile := resolveRequired(defaults.AWSProfile, config.AWSProfile)
		profileBackend := resolveOptional(defaults.AWSProfileBackend, config.AWSProfileBackend)
		profileProvider := resolveOptional(defaults.AWSProfileProvider, config.AWSProfileProvider)
		accountPlan.AWSProfileBackend = resolveRequired(profile, profileBackend)
		accountPlan.AWSProfileProvider = resolveRequired(profile, profileProvider)
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

func resolveOtherAccounts(accounts map[string]config.Account, currentAccount string) map[string]int64 {
	other := make(map[string]int64)
	for name, account := range accounts {
		if name != currentAccount && account.AccountId != nil {
			other[name] = *account.AccountId
		}
	}
	return other
}
