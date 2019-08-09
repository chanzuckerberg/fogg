package examine

import (
	"os"

	"github.com/hashicorp/terraform/config"
)

const apiHostname = "https://registry.terraform.io"
const apiVersion = "/v1"
const resourceType = "/modules/"

//**Local refers to any files located within your local file system**

//GetLocalModules retrieves all terraform modules within a given directory
//TODO:(EC) Define local and global modules OR rename the values
func GetLocalModules(dir string) (*config.Config, error) {
	_, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	config, err := config.LoadDir(dir)
	if err != nil {
		return nil, err
	}

	return config, nil
	// return getAllModules(fs, dir)
}
