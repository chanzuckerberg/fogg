package state

import (
	"fmt"
	"sort"

	"github.com/chanzuckerberg/fogg/config"
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	fogg_hcl "github.com/chanzuckerberg/fogg/exp/hcl"
	"github.com/chanzuckerberg/go-misc/sets"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/lang"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// liberal borrowing from https://github.com/hashicorp/terraform-config-inspect/blob/c481b8bfa41ea9dca417c2a8a98fd21bd0399e14/tfconfig/load_hcl.go#L16
func Run(fs afero.Fs, configFile, path string) error {
	// figure out which component or account we are talking about
	conf, err := config.FindAndReadConfig(fs, configFile)
	if err != nil {
		return err
	}

	component, err := conf.PathToComponentType(path)
	if err != nil {
		return err
	}
	logrus.Debugf("component kind %s", component.Kind)

	// collect remote state references
	references, err := collectRemoteStateReferences(path)
	if err != nil {
		return err
	}
	logrus.Debugf("in %s found references %#v", path, references)

	// for each reference, figure out if it is an account or component, since those are separate in our configs
	accounts := []string{}
	components := []string{}

	// we do accounts for both accounts and env components
	for _, r := range references {
		if _, found := conf.Accounts[r]; found {
			accounts = append(accounts, r)
		}
	}

	if component.Kind == "envs" {
		env := conf.Envs[component.Env]

		for _, r := range references {
			if _, found := env.Components[r]; found {
				components = append(components, r)
			}
		}
	}

	// update fogg.yml with new references
	logrus.Debugf("found accounts %#v", accounts)
	logrus.Debugf("found components %#v", components)

	switch component.Kind {
	case "accounts":
		c := conf.Accounts[component.Name]
		if c.Common.DependsOn == nil {
			c.Common.DependsOn = &v2.DependsOn{}
		}

		c.DependsOn.Accounts = accounts
		conf.Accounts[component.Name] = c
	case "envs":
		c := conf.Envs[component.Env].Components[component.Name]

		if c.Common.DependsOn == nil {
			c.Common.DependsOn = &v2.DependsOn{}
		}

		c.DependsOn.Accounts = accounts
		c.DependsOn.Components = components

		conf.Envs[component.Env].Components[component.Name] = c
	default:
		return fmt.Errorf("unknown component.Kind: %s", component.Kind)
	}

	return conf.Write(fs, configFile)
}

func collectRemoteStateReferences(path string) ([]string, error) {
	references := sets.StringSet{}

	err := fogg_hcl.ForeachBlock(path, func(block *hcl.Block) error {
		logrus.Debugf("block type: %v", block.Type)

		attrs, _ := block.Body.JustAttributes()

		for _, v := range attrs {
			refs, _ := lang.ReferencesInExpr(v.Expr)

			for _, r := range refs {
				if r == nil {
					continue
				}
				logrus.Debugf("ref: %v", r)
				if resource, ok := r.Subject.(addrs.ResourceInstance); ok {
					if resource.Resource.Type == "terraform_remote_state" {
						references.Add(resource.Resource.Name)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	refNames := references.List()

	sort.Strings(refNames)
	return refNames, nil
}
