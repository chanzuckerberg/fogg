package modules

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chanzuckerberg/fogg/config"
	fogg_hcl "github.com/chanzuckerberg/fogg/exp/hcl"
	"github.com/chanzuckerberg/go-misc/sets"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Run will read all terraform files in `path` and return repository-local module references
// The references will be relative to the root of the repository.
// fs is assumed to be rooted at the root of the repository.
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

	// make paths relative to repository root
	absoluteReferences := sets.StringSet{}

	for _, r := range references.List() {
		p := filepath.Join(path, r)
		absoluteReferences.Add(p)
	}

	refNames := absoluteReferences.List()
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

					if strings.HasPrefix(v.AsString(), "./") || strings.HasPrefix(v.AsString(), "../") {
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
