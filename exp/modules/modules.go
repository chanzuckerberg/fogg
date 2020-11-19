package modules

import (
	"fmt"
	"sort"

	"github.com/chanzuckerberg/fogg/config"
	fogg_hcl "github.com/chanzuckerberg/fogg/exp/hcl"
	"github.com/chanzuckerberg/go-misc/sets"
	"github.com/hashicorp/hcl/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

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

	sources, err := collectModuleSources(path)
	if err != nil {
		return err
	}

	logrus.Debugf("found sources %v", sources)

	for _, s := range sources {
		fmt.Println(s)
	}

	return nil
}

func collectModuleSources(path string) ([]string, error) {
	references := sets.StringSet{}

	err := fogg_hcl.ForeachBlock(path, func(block *hcl.Block) error {
		if block.Type == "module" {
			attrs, _ := block.Body.JustAttributes()
			for _, a := range attrs {
				if a.Name == "source" {
					logrus.Debugf("attr %v", a)
					v, err := a.Expr.Value(nil)
					if err != nil {
						return err
					}
					logrus.Debugf("attr expr value %v", v.AsString())
					references.Add(v.AsString())
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
