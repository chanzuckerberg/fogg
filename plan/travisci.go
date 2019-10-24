package plan

import (
	"fmt"
	"sort"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
)

type TravisProject struct {
	Name    string
	Dir     string
	Command string
}
type TravisCI struct {
	Enabled     bool
	FoggVersion string
	TestBuckets [][]TravisProject
	AWSProfiles map[string]AWSRole
	Buildevents bool
}

// TODO(el): mostly a duplicate of buildAtlantis(). refactor later
func (p *Plan) buildTravisCI(c *v2.Config, foggVersion string) TravisCI {
	enabled := false
	buildeventsEnabled := false
	projects := []TravisProject{}
	awsProfiles := map[string]AWSRole{}

	if p.Global.TravisCI.Enabled {
		enabled = true
		buildeventsEnabled = buildeventsEnabled || p.Global.TravisCI.Buildevents

		proj := TravisProject{
			Name:    "global",
			Dir:     "terraform/global",
			Command: p.Global.TravisCI.Command,
		}
		projects = append(projects, proj)

		if p.Global.Backend.AccountID != nil {
			awsProfiles[p.Global.Backend.Profile] = AWSRole{
				AccountID: *p.Global.Backend.AccountID,
				RoleName:  p.Global.TravisCI.AWSRoleName,
			}
			if p.Global.Providers.AWS != nil {
				a := *p.Global.Providers.AWS
				awsProfiles[a.Profile] = AWSRole{
					AccountID: *p.Global.Backend.AccountID,
					RoleName:  p.Global.TravisCI.AWSRoleName,
				}
			}
		}
	}

	for name, acct := range p.Accounts {
		if acct.TravisCI.Enabled {
			enabled = true
			buildeventsEnabled = buildeventsEnabled || p.Global.TravisCI.Buildevents

			proj := TravisProject{
				Name:    fmt.Sprintf("accounts/%s", name),
				Dir:     fmt.Sprintf("terraform/accounts/%s", name),
				Command: acct.TravisCI.Command,
			}

			projects = append(projects, proj)

			// Grab all profiles from accounts
			if acct.Backend.AccountID != nil {
				awsProfiles[acct.Backend.Profile] = AWSRole{
					AccountID: *acct.Backend.AccountID,
					RoleName:  acct.TravisCI.AWSRoleName,
				}
			}
			if acct.Providers.AWS != nil {
				awsProfiles[acct.Providers.AWS.Profile] = AWSRole{
					AccountID: acct.Providers.AWS.AccountID.String(),
					RoleName:  acct.TravisCI.AWSRoleName,
				}
			}
		}
	}

	for _, env := range p.Envs {
		for _, c := range env.Components {
			if c.TravisCI.Enabled {
				enabled = true
				buildeventsEnabled = buildeventsEnabled || p.Global.TravisCI.Buildevents

				// proj := TravisProject{
				// 	Name:    fmt.Sprintf("%s/%s", envName, cName),
				// 	Dir:     fmt.Sprintf("terraform/envs/%s/%s", envName, cName),
				// 	Command: "lint",
				// }

				// projects = append(projects, proj)

				if c.Backend.AccountID != nil {
					awsProfiles[c.Backend.Profile] = AWSRole{
						AccountID: *c.Backend.AccountID,
						RoleName:  c.TravisCI.AWSRoleName,
					}
				}

				if c.Providers.AWS != nil {
					a := *c.Providers.AWS
					awsProfiles[a.Profile] = AWSRole{
						AccountID: a.AccountID.String(),
						RoleName:  c.TravisCI.AWSRoleName,
					}
				}
			}
		}
	}

	// for moduleName := range p.Modules {
	// 	proj := TravisProject{
	// 		Name:    fmt.Sprintf("modules/%s", moduleName),
	// 		Dir:     fmt.Sprintf("terraform/modules/%s", moduleName),
	// 		Command: "lint",
	// 	}
	// 	projects = append(projects, proj)
	// }

	var buckets int
	if c.Defaults.Tools != nil &&
		c.Defaults.Tools.TravisCI != nil &&
		c.Defaults.Tools.TravisCI.TestBuckets != nil &&
		*c.Defaults.Tools.TravisCI.TestBuckets > 0 {

		buckets = *c.Defaults.Tools.TravisCI.TestBuckets
	} else {
		buckets = 1
	}

	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})

	testBuckets := make([][]TravisProject, buckets)
	for i, proj := range projects {
		bucket := i % buckets
		testBuckets[bucket] = append(testBuckets[bucket], proj)
	}

	tr := TravisCI{
		Enabled:     enabled,
		Buildevents: buildeventsEnabled,
		FoggVersion: foggVersion,
		TestBuckets: testBuckets,
		AWSProfiles: awsProfiles,
	}
	return tr
}
