package versioning

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/afero"

	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

//	"gopkg.in/src-d/go-git.v4"

/*
 * 1) Read main file
 * 2) Find the source for each module
 *   a) If that file is a directory read that directory
 *   b) Otherwise, check if it is github
 * 3) Store the contents of each module
 */

/*
 *SOLVED: Path strings will get extremely long as submodules are parsed (can be fixed)
 *Solution -> Use Afero
 *SOLVED: The chanzuckerberg cztack links don't seem to link to anything (can be fixed, but getting the latest version will be challenging)
 *Solution -> Go Getter
 *Problem 3: Local files do not have versions
 *SOLVED: Terraform registry files do not always say they're from the registry
 *Solution -> If the file is not recognized, then parse fogg.yml and see if they are module sources
 *Problem 4: Do I want to get
 */

type ModuleWrapper struct {
	moduleSource string
	version      string
	module       *tfconfig.Module
}

//GetLocalModules Retrieves all modules that the given directory depends on
func GetLocalModules(path string) []ModuleWrapper {
	modules, err := retrieveAllDependencies(path)
	if err != nil {
		panic(err)
	}
	return modules
}

//retrieveAllDependencies retrieves all modules and resources used to assemble the root module
func retrieveAllDependencies(path string) ([]ModuleWrapper, error) {
	var temp []ModuleWrapper

	//Read the initial path and add the submodules to directories
	modules, err := findSubmodules(path)
	if err != nil {
		return nil, err
	}
	temp = append(temp, modules...)

	for { //Do: add modules to modules slice, While: temp is not empty

		mods, err := findSubmodules(temp[0].module.Path)
		if err == nil {
			//Add new elements to temp
			//TODO: Dont add repeats
			temp = append(temp, mods...)
			modules = append(modules, mods...)
		}
		//remove the current element
		temp = append(temp[:0], temp[1:]...)

		if !(len(temp) > 0) {
			break
		}
	}

	return modules, nil
}

//findSubmodules retrieves all modules related to the given path
func findSubmodules(dir string) ([]ModuleWrapper, error) {
	_, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	//FIXME: Make contains check module path
	// if(modules != nil && contains(modules, dir))
	var modules []ModuleWrapper
	mod, diag := tfconfig.LoadModule(dir)
	if diag.HasErrors() {
		return nil, err
	}

	//Loads all modules that the current module depends on
	keys := getModules(mod)
	sources := GetSources(mod, keys)
	for i, source := range sources {
		// If it is a github link, do not append a local directory
		if strings.HasPrefix(source.moduleSource, "github.com") {
			//FIXME: I am doing all github specific actions here, should I do them somewhere else?
			mod, err := getFromGithub(source.moduleSource)
			if err != nil {
				fmt.Println(err)
				continue
			}

			sources[i].module = mod
			sources[i].version = getVersion(sources[i])
			modules = append(modules, sources[i])
		} else if strings.HasPrefix(source.moduleSource, "terraform-aws-modules/") {
			mod, err := downloadModule(source.moduleSource, source.version)
			if err != nil {
				fmt.Println(err)
				continue
			}

			sources[i].module = mod
			modules = append(modules, sources[i])
		} else {
			//Append the local directory to the file
			//FIXME: Appends to an index that does not exist
			mod, diag := tfconfig.LoadModule(dir + source.moduleSource + "/")
			if diag.HasErrors() {
				fmt.Println(err)
				continue
			}

			sources[i].module = mod
			modules = append(modules, sources[i])
		}
	}
	return modules, nil
}

func GetSources(mod *tfconfig.Module, keys []string) []ModuleWrapper {
	modMap := mod.ModuleCalls

	//Make 3 functions that separate git and other sources
	sources := make([]ModuleWrapper, 0)
	for _, key := range keys {
		//FIXME: Fix any modules that might not contain verison
		sources = append(sources, ModuleWrapper{modMap[key].Source, modMap[key].Version, nil})
	}
	return sources
}

func getModules(module *tfconfig.Module) []string {
	modMap := module.ModuleCalls
	keys := make([]string, 0)
	for key := range modMap {
		keys = append(keys, key)
	}

	return keys
}

//GetFromGithub Retrieves modules that are available through github
func getFromGithub(repo string) (*tfconfig.Module, error) {
	pwd, e := os.Getwd()
	if e != nil {
		return nil, e
	}

	//TODO: Pass the fs from the top level
	//Create temporary directory
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
	//TODO: Make dir name, module repo name
	tmpDirPath, err := afero.TempDir(fs, ".", "cztack")
	if err != nil {
		fmt.Println("There was an error creating tmp Dir")
		return nil, err
	}
	defer os.RemoveAll(tmpDirPath)

	//Load github file
	err = getter.Get(tmpDirPath, repo)
	if err != nil {
		fmt.Println("There was an issue getting the repo")
		return nil, err
	}

	//Read the files into a directory
	files, err := afero.ReadDir(fs, tmpDirPath)
	if err != nil {
		fmt.Println("There was an issue reading the directory")
		return nil, err
	}
	fmt.Println(files)

	//Read the module
	//FIXME: Return value, potentially see whats in Diagnostics
	mod, diag := tfconfig.LoadModule(tmpDirPath)
	if diag.HasErrors() {
		return nil, nil
	}
	//TODO: Returns diagnostics error
	if err != nil {
		fmt.Println("tconfig could not read tmpDir")
		return nil, err
	}

	return mod, nil
}

//downloadModule retrieves terraform modules from the registry
func downloadModule(modulePath string, version string) (*tfconfig.Module, error) {
	baseUrl := "https://registry.terraform.io/v1/modules/"
	pwd, e := os.Getwd()
	if e != nil {
		return nil, e
	}

	//TODO: Pass the fs from the top level
	//Create temporary directory
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
	//TODO: Make dir name, module repo name
	tmpDirPath, err := afero.TempDir(fs, ".", "registry-download")
	if err != nil {
		fmt.Println("There was an error creating tmp Dir")
		return nil, err
	}
	defer os.RemoveAll(tmpDirPath)

	//Load github file
	err = getter.Get(tmpDirPath, baseUrl+modulePath+"/"+version+"/download")
	if err != nil {
		fmt.Println("There was an issue getting the repo")
		return nil, err
	}

	//Read the files into a directory
	files, err := afero.ReadDir(fs, tmpDirPath)
	if err != nil {
		fmt.Println("There was an issue reading the directory")
		return nil, err
	}
	fmt.Println(files)

	//Read the module
	//FIXME: Return value, potentially see whats in Diagnostics
	mod, diag := tfconfig.LoadModule(tmpDirPath)
	if diag.HasErrors() {
		return nil, nil
	}
	//TODO: Returns diagnostics error
	if err != nil {
		fmt.Println("tconfig could not read tmpDir")
		return nil, err
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

//TODO: Potentially do something when empty string is returned
func getVersion(mod ModuleWrapper) string {
	return strings.Split(mod.moduleSource, "ref=v")[1]
}
