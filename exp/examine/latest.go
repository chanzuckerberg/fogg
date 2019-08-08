package examine

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

//LatestModuleVersions retrieves the latest version of the provided modules
func LatestModuleVersions(fs afero.Fs, modules []ModuleWrapper) ([]ModuleWrapper, error) {
	var latestModules []ModuleWrapper
	var module ModuleWrapper
	var resource string
	var err error

	for _, mod := range modules {
		if strings.HasPrefix(mod.moduleSource, githubURL) { //If the resource is from github then retrieve from custom url
			resource, err = generateURL(fs, mod)
			if err != nil {
				return nil, errs.WrapUserf(err, "Could not generate url for %s latest module", mod.moduleSource)
			}
			module.module, err = GetFromGithub(fs, resource)
			if err != nil {
				return nil, err
			}

			latestModules = append(latestModules, module)
		} else if strings.HasPrefix(module.moduleSource, awsRegistry) && false { //DISABLED
			//TODO:(EC) Enable
			resource = module.moduleSource
			res, err := http.Get(resource)
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}

			//FIXME: unmarshal body into ModuleWrapper, specifically in the RegistryModule part
			//TODO:(EC) Get input about adding RegistryModule as a part of module
			err = json.Unmarshal(body, &module)
			if err != nil {
				return nil, err
			}
			latestModules = append(latestModules, module)
		}
	}

	return latestModules, nil
}

//generateURL creates github url for the given module
func generateURL(fs afero.Fs, module ModuleWrapper) (string, error) {
	url := ""
	if strings.HasPrefix(module.moduleSource, "github.com/chanzuckerberg") { //Clone the repo and get the latest tag
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
		url = strings.Split(module.moduleSource, tagPattern)[0] + tagPattern + strings.Split(url, "tags/v")[1]
	}
	//TODO:(EC) Create process for other link types
	return url, nil
}
