package versioning

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/spf13/afero"
)

type ModuleWrapper struct {
	moduleSource string
	version      string
	module       *tfconfig.Module
}

//FIXME: Make struct name more descriptive
type RegistryModule struct {
	ID         string      `json:"id"`
	Namespace  string      `json:"namespace"`
	Name       string      `json:"name"`
	Version    string      `json:"version"`
	Provider   string      `json:"provider"`
	Source     string      `json:"source"`
	Tag        string      `json:"tag"`
	Root       Root        `json:"root"`
	Submodules []Submodule `json:"submodules"`
	Providers  []string    `json:"providers"`
	Versions   []string    `json:"versions"`
}

type Root struct {
	Inputs    []Input    `json:"inputs"`
	Outputs   []Output   `json:"outputs"`
	Resources []Resource `json:"resources"`
}

type Input struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type Output struct {
	Name string `json:"name"`
}

type Resource struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Submodule struct {
	Path      string     `json:"path"`
	Name      string     `json:"name"`
	Inputs    []Input    `json:"inputs"`
	Outputs   []Output   `json:"outputs"`
	Resources []Resource `json:"resources"`
}

const githubURL = "github.com"
const awsRegistry = "terraform-aws-modules/"
const cztack = "/chanzuckerberg/cztack/"
const protocol = "https://"

//GetFromGithub Retrieves modules that are available through github
func GetFromGithub(fs afero.Fs, repo string) (*tfconfig.Module, error) {

	//Create temporary directory
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
