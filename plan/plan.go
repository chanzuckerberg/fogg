package plan

import (
	"github.com/ryanking/fogg/config"
	"github.com/ryanking/fogg/util"
	"github.com/spf13/afero"
)

type common struct {
	AWSRegion          string
	InfraBucket        string
	Project            string
	SharedInfraPath    string
	TerraformVersion   string
	AWSRegions         string
	SharedInfraVersion string
	Owner              string
	AWSProfileBackend  string
	AWSProfileProvider string
}

type account struct {
	AccountName   string
	OtherAccounts string

	common
}

type plan struct {
	accounts map[string]account
}

func Plan(fs afero.Fs) (*plan, error) {
	c, _ := config.FindAndReadConfig(fs)
	util.Dump(c)
	// read config and validate
	// build repo plan
	// build .sicc version plan
	buildAccounts(c)
	// build modules plan
	// build envs plan
	// walk config and apply inheritance rules
	return nil, nil
}

func buildAccounts(c *config.Config) (map[string]*account, error) {
	accountPlans := make(map[string]*account, len(c.Accounts))
	for name, config := range c.Accounts {
		util.Dump(name)
		util.Dump(config)
		accountPlan := &account{}
		// copy defaults
		// fix shared infra base
		// resolve account name
		// resolve profiles
		// resolve other accounts
		accountPlans[name] = accountPlan
	}

	return accountPlans, nil
}

func coalesceStrings(in []*string) *string {
	for i := 0; i < len(in); i++ {
		if in[i] != nil {
			return in[i]
		}
	}
	return nil
}
