package examine

import (
	"os"
	"strings"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/hashicorp/terraform/config"
	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

//LatestModuleVersions retrieves the latest version of the provided modules
func LatestModuleVersions(fs afero.Fs, config *config.Config) ([]ModuleWrapper, error) {
	var latestModules []ModuleWrapper
	var module ModuleWrapper
	var resource string
	var err error

	for _, mod := range config.Modules {
		if strings.HasPrefix(mod.Source, githubURL) { //If the resource is from github then retrieve from custom url
			resource, err = generateURL(fs, mod)
			if err != nil {
				return nil, errs.WrapUserf(err, "Could not generate url for %s latest module", mod.Source)
			}
			module.module, err = GetFromGithub(fs, resource)
			if err != nil {
				return nil, err
			}
			module.moduleSource = mod.Source
			module.version = mod.Version

			latestModules = append(latestModules, module)
		}
		//TODO: Implement retrieving the latest module from the terraform registry
	}

	return latestModules, nil
}

//generateURL creates github url for the given module
func generateURL(fs afero.Fs, module *config.Module) (string, error) {
	url := ""
	if strings.HasPrefix(module.Source, "github.com/chanzuckerberg") { //Clone the repo and get the latest tag
		tmpDir, err := afero.TempDir(fs, ".", "github")
		if err != nil {
			return "", err
		}
		defer os.RemoveAll(tmpDir)

		repo, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
			URL: protocol + githubURL + cztack,
		})
		if err != nil {
			return "", errs.WrapUser(err, "Could not clone repo")
		}

		//Runs git ls-remote --tags
		tagIterator, err := repo.Tags()
		if err != nil {
			return "", errs.WrapUser(err, "Could not find tags for repo")
		}

		err = tagIterator.ForEach(func(t *plumbing.Reference) error {
			url = t.Name().String()
			return nil
		})
		//TODO:(EC) Make the tag naming scheme modular
		url = strings.Split(module.Source, tagPattern)[0] + tagPattern + strings.Split(url, "tags/v")[1]
	}
	//TODO:(EC) Create process for other link types
	return url, nil
}
