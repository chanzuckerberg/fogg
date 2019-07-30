package versioning

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

//	"github.com/hashicorp/terraform-config-inspect/tfconfig"

const registry = "https://registry.terraform.io/v1/modules"
const resource = "https://registry.terraform.io/v1/modules/terraform-aws-modules/alb/aws"

func GetModule() Module {
	var module Module

	res, err := http.Get(resource)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &module)
	if err != nil {
		panic(err)
	}

	return module
}

func GetAWSModules() {
	//https://registry.terraform.io/v1/modules?namespace=terraform-aws-modules
	//
	//https://registry.terraform.io/v1/modules?namespace=terraform-aws-modules&offset=15
	//https://registry.terraform.io/v1/modules?provider=aws&verified=true
	res, err := http.Get("https://registry.terraform.io/v1/modules?namespace=terraform-aws-modules")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	bytes, _ := ioutil.ReadAll(res.Body)
	// str := string(bytes)

	var j interface{}

	e := json.Unmarshal(bytes, &j)
	if e != nil {
	}

	fmt.Println(j)
}

// GetGlobalModules Retrieves modules related to tconfig from the registry
func GetGlobalModules(modules []ModuleWrapper) []ModuleWrapper {
	var globalModules []ModuleWrapper
	var module ModuleWrapper

	for _, mod := range modules {
		//FIXME: Make this work for more than
		if !strings.HasPrefix(mod.moduleSource, "github.com") {
			continue
		}
		resource := generateUrl(mod)
		if strings.HasPrefix(mod.moduleSource, "github.com") {
			module.module, _ = GetFromGithub(resource)

		} else {
			res, err := http.Get(resource)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}

			//FIXME: unmarshal into ModuleWrapper Module struct
			err = json.Unmarshal(body, &module)
			if err != nil {
				panic(err)
			}
		}

		globalModules = append(globalModules, module)
	}

	return globalModules
}

//Takes the link
func generateUrl(module ModuleWrapper) string {
	url := ""
	if strings.HasPrefix(module.moduleSource, "github.com/chanzuckerberg") {
		//Take the link
		pwd, e := os.Getwd()
		if e != nil {
			return ""
		}

		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
		tmpDir, err := afero.TempDir(fs, ".", "github")
		if err != nil {
			fmt.Println("There was an error creating tmp Dir")
			return ""
		}
		defer os.RemoveAll(tmpDir)

		repo, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
			URL: "https://github.com/chanzuckerberg/cztack/",
		})
		if err != nil {
			fmt.Println("Could not clone repo.")
			return ""
		}

		//Run git ls-remote --tags
		tagIterator, err := repo.Tags()
		if err != nil {
			fmt.Println("Could not find tags")
			return ""
		}

		// fmt.Printf("Calling next returns: %v", tagIterator.
		err = tagIterator.ForEach(func(t *plumbing.Reference) error {
			url = t.Name().String()
			return nil
		})

		url = strings.Split(module.moduleSource, "ref=v")[0] + "ref=v" + strings.Split(url, "tags/v")[1]
	}
	//TODO: For other link types include https://
	fmt.Println(url)
	return url
}
