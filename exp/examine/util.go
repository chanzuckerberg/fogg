package examine

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/spf13/afero"
)

// TODO:(EC) Add a RegistryModule field
type ModuleWrapper struct {
	moduleSource string
	version      string
	module       *tfconfig.Module
}

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
const tagPattern = "ref="

// GetFromGithub Retrieves modules that are available through github
func GetFromGithub(fs afero.Fs, repo string) (*tfconfig.Module, error) {
	//FIXME: (EC) Create temporary directory, when tests fail directory stays
	//TODO: Make directory name more general
	tmpDir, err := afero.TempDir(fs, ".", "cztack")
	if err != nil {
		return nil, errs.WrapInternal(err, "There was an issue creating a temp directory")
	}
	defer os.RemoveAll(tmpDir)

	//Load github file
	err = getter.Get(tmpDir, repo)
	if err != nil {
		return nil, errs.WrapInternal(err, "There was an issue getting the repo")
	}

	//Load module into struct
	mod, diag := tfconfig.LoadModule(tmpDir)
	if diag.HasErrors() {
		return nil, errs.WrapInternal(diag.Err(), "There was an issue loading the module")
	}

	return mod, nil
}
