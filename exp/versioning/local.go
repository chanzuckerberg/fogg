package versioning

import (
	"fmt"
	"os"
	"strings"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/spf13/afero"
)

const apiHostname = "https://registry.terraform.io"
const apiVersion = "/v1/modules/"

//TODO: Replace some of the strings with modular constants
//TODO: Freeze the code by creating a test directory

//Local refers to any files located within your local file system.

//GetLocalModules retrieves all terraform modules within a local directory
//TODO:(EC) Define local and global modules OR rename the values
func GetLocalModules(fs afero.Fs, dir string) ([]ModuleWrapper, error) {
	_, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	modules, err := getAllModules(fs, dir)
	if err != nil {
		return nil, err
	}
	return modules, nil
}

//getAllDependencies retrieves all modules from the root directory
func getAllModules(fs afero.Fs, dir string) ([]ModuleWrapper, error) {
	var queue []ModuleWrapper

	//Find submodules for the root directory
	submodules, err := getSubmodules(fs, dir)
	if err != nil {
		return nil, err
	}
	queue = append(queue, submodules...)

	for { //Do: recursively add new submodules to queue, While: queue is not empty

		submods, err := getSubmodules(fs, queue[0].module.Path)
		if err == nil {
			//Add new elements to temp
			//FIXME:(EC) Dont add repeats
			queue = append(queue, submods...)
			submodules = append(submodules, submods...)
		}
		//remove the current element
		queue = append(queue[:0], queue[1:]...)

		if !(len(queue) > 0) {
			break
		}
	}

	return submodules, nil
}

//getSubmodules retrieves all submodules related to the root directory
func getSubmodules(fs afero.Fs, dir string) ([]ModuleWrapper, error) {
	_, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	//FIXME: Make contains check module path
	// if(modules != nil && contains(modules, dir))
	var modules []ModuleWrapper
	mod, diag := tfconfig.LoadModule(dir)
	if diag.HasErrors() {
		fmt.Println(diag.Err())
		return nil, diag.Err()
	}

	//Loads all modules that the current module depends on
	sources := getSources(mod)
	for i, source := range sources {
		if strings.HasPrefix(source.moduleSource, githubURL) { // If the module is a github link, specifically cztack
			mod, err := GetFromGithub(fs, source.moduleSource)
			if err != nil {
				return nil, err
			}

			sources[i].module = mod
			sources[i].version = getVersion(sources[i])
			modules = append(modules, sources[i])
		} else if strings.HasPrefix(source.moduleSource, awsRegistry) { // If the module is from the tf registry
			mod, err := downloadModule(fs, source.moduleSource, source.version)
			if err != nil {
				return nil, err
			}

			sources[i].module = mod
			modules = append(modules, sources[i])
		} else { //Otherwise, the module is not the leaf
			//Append the local directory to the file
			mod, diag := tfconfig.LoadModule(dir + source.moduleSource + "/")
			if diag.HasErrors() {
				return nil, diag.Err()
			}

			sources[i].module = mod
			modules = append(modules, sources[i])
		}
	}
	return modules, nil
}

//getSources retrieves the source for each module and creates a ModuleWrapper with it
func getSources(mod *tfconfig.Module) []ModuleWrapper {
	//ModuleCalls represents a module's submodules
	modMap := mod.ModuleCalls
	sources := make([]ModuleWrapper, 0)

	//TODO:(EC) Make 3 functions that separate git and other sources
	for key := range modMap {
		//FIXME:(EC) Fix any modules that might not contain verison
		sources = append(sources, ModuleWrapper{modMap[key].Source, modMap[key].Version, nil})
	}

	return sources
}

//downloadModule downloads terraform module data from the tf registry
func downloadModule(fs afero.Fs, modulePath string, version string) (*tfconfig.Module, error) {
	tmpDir, err := afero.TempDir(fs, ".", "registry-download")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	//download the module from tf registry
	err = getter.Get(tmpDir, apiHostname+apiVersion+modulePath+"/"+version+"/download")
	if err != nil {
		return nil, errs.WrapUser(err, "There was an issue downloading the file")
	}

	mod, diag := tfconfig.LoadModule(tmpDir)
	if diag.HasErrors() {
		return nil, errs.WrapUser(diag.Err(), "Could not load module")
	}

	return mod, nil
}

func contains(modules []*tfconfig.Module, mod *tfconfig.Module) bool {
	//FIXME: Make sure contains does not just check path but module name
	for _, m := range modules {
		if m == mod {
			return true
		}
	}
	return false
}

//TODO:(EC) Error handle when empty string is encountered
func getVersion(mod ModuleWrapper) string {
	return strings.Split(mod.moduleSource, "ref=v")[1]
}
