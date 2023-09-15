package config

import (
	"encoding/json"
	"io"
	"path/filepath"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/util"
	goVersion "github.com/hashicorp/go-version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v3"
)

var defaultTerraformVersion = goVersion.Must(goVersion.NewVersion("1.2.6"))

// DefaultFoggVersion is the version that fogg will generate by default
const DefaultFoggVersion = 2

func defaultEnabled(a bool) *bool {
	return &a
}

// InitConfig initializes the config file using user input
func InitConfig(project, region, bucket, table, awsProfile, owner, awsAccountID *string, awsProviderVersion string) *v2.Config {
	accountID := json.Number(*awsAccountID)

	return &v2.Config{
		Defaults: v2.Defaults{
			Common: v2.Common{
				Backend: &v2.Backend{
					Bucket:      bucket,
					Profile:     awsProfile,
					Region:      region,
					DynamoTable: table,
				},
				Owner:   owner,
				Project: project,
				Providers: &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID: &accountID,
						Profile:   awsProfile,
						Region:    region,
						CommonProvider: v2.CommonProvider{
							Enabled: defaultEnabled(true),
							Version: &awsProviderVersion,
						},
					},
				},
				TerraformVersion: util.StrPtr(defaultTerraformVersion.String()),
			},
		},
		Accounts: map[string]v2.Account{},
		Envs:     map[string]v2.Env{},
		Modules:  map[string]v2.Module{},
		Version:  DefaultFoggVersion,
	}
}

// FindConfig loads a config and its version into memory
func FindConfig(fs afero.Fs, configFile string) ([]byte, int, error) {
	f, err := fs.Open(configFile)
	if err != nil {
		return nil, 0, errs.NewUserf("could not open %s config file", configFile)
	}
	defer f.Close()

	b, e := io.ReadAll(f)
	if e != nil {
		return nil, 0, errs.WrapUser(e, "unable to read config")
	}

	v, err := detectVersion(b, fs, configFile)
	if err != nil {
		return nil, 0, err
	}
	logrus.Debugf("config file version: %#v\n", v)
	return b, v, nil
}

// FindAndReadConfig locates config file and reads it based on the version
func FindAndReadConfig(fs afero.Fs, configFile string) (*v2.Config, error) {
	b, v, err := FindConfig(fs, configFile)
	if err != nil {
		return nil, err
	}

	switch v {
	case 2:
		return v2.ReadConfig(fs, b, configFile)
	default:
		return nil, errs.NewUser("could not figure out config file version")
	}
}

type ver struct {
	Version int `yaml:"version"`
}

// detectVersion will detect the version of a config
func detectVersion(b []byte, fs afero.Fs, configFile string) (int, error) {
	v := &ver{}
	var err error

	info, err := fs.Stat(configFile)
	if err != nil {
		return 0, errs.WrapUserf(err, "unable to stat %s", configFile)
	}

	//Unmarshals based on file extension
	switch filepath.Ext(info.Name()) {
	case ".yml", ".yaml":
		err = yaml.Unmarshal(b, v)
	default:
		return 0, errs.NewUser("File type is not supported")
	}

	if err != nil {
		return 0, err
	}
	if v.Version == 0 {
		return 2, nil
	}
	return v.Version, nil
}
