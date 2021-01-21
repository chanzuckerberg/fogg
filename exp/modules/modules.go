package modules

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chanzuckerberg/fogg/config"
	fogg_hcl "github.com/chanzuckerberg/fogg/exp/hcl"
	"github.com/chanzuckerberg/go-misc/sets"
	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/hcl/v2"
	"github.com/pkg/errors"
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
	references, err := extractLocalModules(path)
	if err != nil {
		return nil, err
	}
	refNames := references.List()
	sort.Strings(refNames)
	return refNames, nil
}

// extractLocalModules will recursively look for all local module references
// we assume there are no cycles in this dependency graph
func extractLocalModules(path string) (sets.StringSet, error) {
	references := sets.StringSet{}

	// for a given path, find all module references
	err := fogg_hcl.ForeachBlock(path, func(block *hcl.Block) error {
		if block.Type == "module" {
			attrs, _ := block.Body.JustAttributes()
			for _, a := range attrs {
				if a.Name == "source" {
					v, err := a.Expr.Value(nil)
					if err != nil {
						return err
					}

					// we only want local modules, so use go-getter to detect "file:" type
					str, err2 := getter.Detect(v.AsString(), path, getter.Detectors)
					if err2 != nil {
						return errors.Wrap(err2, "unable to detect module type")
					}
					if strings.HasPrefix(str, "file:") {
						references.Add(v.AsString())
					}
				}
			}
		}
		return nil
	})

	// recursively find modules
	for _, r := range references.List() {
		t := filepath.Join(path, r)
		refs, err := extractLocalModules(t)
		if err != nil {
			return references, err
		}

		// convert paths to be relative to the original path for use in tfe trigger prefixes
		for _, x := range refs.List() {
			p, _ := filepath.Rel(path, filepath.Join(t, x))
			references.Add(p)
		}
	}
	return references, err
}
