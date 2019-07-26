package versioning 

import(
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	
)
//	"github.com/hashicorp/terraform-config-inspect/tfconfig"

const registry = "https://registry.terraform.io/v1/modules"
const resource = "https://registry.terraform.io/v1/modules/terraform-aws-modules/alb/aws"

func GetModule() Module{
	var module Module

	res, err := http.Get(resource)
	if err != nil{
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil{
		panic(err)
	}

	err = json.Unmarshal(body, &module)
	if err != nil{
		panic(err)
	}

	return module
}

func GetAWSModules(){
	//https://registry.terraform.io/v1/modules?namespace=terraform-aws-modules
//
//https://registry.terraform.io/v1/modules?namespace=terraform-aws-modules&offset=15
//https://registry.terraform.io/v1/modules?provider=aws&verified=true
	res, err := http.Get("https://registry.terraform.io/v1/modules?namespace=terraform-aws-modules")
	if err != nil{
		panic(err)
	}
	defer res.Body.Close()

	bytes, _ := ioutil.ReadAll(res.Body)
	// str := string(bytes)

	var j interface{}

	e := json.Unmarshal(bytes, &j)
	if e != nil{}

	fmt.Println(j)
}

//GetGlobalModules Retrieves modules related to tconfig from the registry
// func GetGlobalModules(modules []*tfconfig.Module) []Module{
// 	var globalModules []Module
// 	var module Module

// 	for _, mod := range modules{
		
// 		// resource := createHttp(mod)
// 		res, err := http.Get(resource)
// 		if err != nil{
// 			panic(err)
// 		}
// 		defer res.Body.Close()

// 		body, err := ioutil.ReadAll(res.Body)
// 		if err != nil{
// 			panic(err)
// 		}

// 		err = json.Unmarshal(body, &module)
// 		if err != nil{
// 			panic(err)
// 		}

// 		globalModules = append(globalModules, module)
// 	}
	

// 	return globalModules
// }


//
// func createHttp(module *tfconfig.Module){
// 	resource := "/terraform-aws-modules"
// 	var middleWare string
// 	provider := "/aws"
// }