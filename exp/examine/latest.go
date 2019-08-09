package examine

import (
	"context"
	"strings"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/config"
	"github.com/spf13/afero"
)

//LatestModuleVersions retrieves the latest version of the provided modules
func LatestModuleVersions(fs afero.Fs, config *config.Config) ([]ModuleWrapper, error) {
	var latestModules []ModuleWrapper
	var module ModuleWrapper
	var resource string
	var err error

	for _, mod := range config.Modules {
		if strings.HasPrefix(mod.Source, githubURL) { //If the resource is from github then retrieve from custom url
			resource, err = createGitUrl(mod)
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

//createGitUrl retrieves the latest release version and creates an HTTP accessible link
func createGitUrl(module *config.Module) (string, error) {
	splitString := strings.Split(module.Source, "/")
	owner, repo := splitString[1], splitString[2]

	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		return "", err
	}

	return strings.Split(module.Source, tagPattern)[0] + tagPattern + *release.TagName, nil
}
