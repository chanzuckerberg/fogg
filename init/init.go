package init

import (
	"github.com/chanzuckerberg/fogg/config"
	prompt "github.com/segmentio/go-prompt"
	"github.com/spf13/afero"
)

const AWSProviderVersion = "2.47.0"

func userPrompt() (string, string, string, string, string, string) {
	project := prompt.StringRequired("project name?")
	region := prompt.StringRequired("aws region?")
	bucket := prompt.StringRequired("infra bucket name?")
	table := prompt.String("infra dynamo table name?")
	profile := prompt.StringRequired("auth profile?")
	owner := prompt.StringRequired("owner?")

	return project, region, bucket, table, profile, owner
}

//Init reads user console input and generates a fogg.yml file
func Init(fs afero.Fs) error {
	project, region, bucket, table, profile, owner := userPrompt()
	config := config.InitConfig(project, region, bucket, table, profile, owner, AWSProviderVersion)
	e := config.Write(fs, "fogg.yml")
	return e
}
