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

//GetFromGithub Retrieves modules that are available through github
func GetFromGithub(repo string) (*tfconfig.Module, error) {
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
