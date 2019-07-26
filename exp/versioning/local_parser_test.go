package versioning

import(
	"fmt"
	"testing"

)

func TestGetLocalModules(t *testing.T){
	modules := GetLocalModules(theWorks2)

	for _,module := range modules{
		fmt.Printf("\nHere is the module: %v\n", module)
		fmt.Printf("Path: %v\n\n", module.module.Path)
	}
}

func TestGetFromGithub(t *testing.T){
	mod,  err := getFromGithub(repo)
	fmt.Println(err)
	fmt.Println(mod)
}

func TestGetRegistryFromGithub(t *testing.T){
	mod, err := getFromGithub(repo)
	fmt.Println(err)
	fmt.Println(mod)
}

func TestDownloadRegistryFromGithub(t *testing.T){
	path := "terraform-aws-modules/security-group/aws"
	mod, err := downloadModule(path, "2.9.0")
	fmt.Println(err)
	fmt.Println(mod)
}

func TestGetFromAlbHttp(t *testing.T){
	mods, err := retrieveAllDependencies(path)
	if mods == nil || err == nil{}
}
