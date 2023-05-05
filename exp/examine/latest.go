package examine

import (
	"context"
	"strings"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/google/go-github/v27/github"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/spf13/afero"
)

// LatestModuleVersions retrieves the latest version of the provided modules
func LatestModuleVersions(fs afero.Fs, module *tfconfig.Module) ([]ModuleWrapper, error) {
	var latestModules []ModuleWrapper
	var moduleWrapper ModuleWrapper
	var resource string
	var err error

	for _, moduleCall := range module.ModuleCalls {
		if strings.HasPrefix(moduleCall.Source, githubURL) { //If the resource is from github then retrieve from custom url
			resource, err = createGitURL(moduleCall)
			if err != nil {
				return nil, errs.WrapUserf(err, "Could not generate url for %s latest module", moduleCall.Source)
			}
			moduleWrapper.module, err = GetFromGithub(fs, resource)
			if err != nil {
				return nil, err
			}
			moduleWrapper.moduleSource = moduleCall.Source
			moduleWrapper.version = moduleCall.Version

			latestModules = append(latestModules, moduleWrapper)
		}
		//TODO: Implement retrieving the latest module from the terraform registry
	}

	return latestModules, nil
}

// createGitURL retrieves the latest release version and creates an HTTP accessible link
func createGitURL(moduleCall *tfconfig.ModuleCall) (string, error) {
	splitString := strings.Split(moduleCall.Source, "/")
	owner, repo := splitString[1], splitString[2]

	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		return "", err
	}

	return strings.Split(moduleCall.Source, tagPattern)[0] + tagPattern + *release.TagName, nil
}
