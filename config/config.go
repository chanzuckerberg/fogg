package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	v1 "github.com/chanzuckerberg/fogg/config/v1"
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

func InitConfig(project, region, bucket, table, awsProfile, owner, awsProviderVersion string) *v1.Config {
	return &v1.Config{
		Defaults: v1.Defaults{
			AWSProfileBackend:  awsProfile,
			AWSProfileProvider: awsProfile,
			AWSProviderVersion: awsProviderVersion,
			AWSRegionBackend:   region,
			AWSRegionProvider:  region,
			ExtraVars:          map[string]string{},
			InfraBucket:        bucket,
			InfraDynamoTable:   table,
			Owner:              owner,
			Project:            project,
			TerraformVersion:   "0.11.7",
		},
		Accounts: map[string]v1.Account{},
		Docker:   false,
		Envs:     map[string]v1.Env{},
		Modules:  map[string]v1.Module{},
	}
}

// FindConfig loads a config and its version into memory
func FindConfig(fs afero.Fs, configFile string) ([]byte, int, error) {
	f, err := fs.Open(configFile)
	if err != nil {
		return nil, 0, errs.WrapUser(err, "unable to open config file")
	}
	defer f.Close()

	b, e := ioutil.ReadAll(f)

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

func FindAndReadConfig(fs afero.Fs, configFile string) (*v2.Config, error) {
	b, v, err := FindConfig(fs, configFile)
	if err != nil {
		return nil, err
	}

	switch v {
	case 1: //Upgrade the config version

		c, err := v1.ReadConfig(b)
		if err != nil {
			return nil, err
		}

		return UpgradeConfigVersion(c)
	case 2:

		return v2.ReadConfig(fs, b, configFile)
	default:
		return nil, errs.NewUser("could not figure out config file version")
	}

}

type ver struct {
	Version int `json:"version" yaml:"version"`
}

// detectVersion will detect the version of a config
func detectVersion(b []byte, fs afero.Fs, configFile string) (int, error) {
	v := &ver{}
	var err error

	info, err := fs.Stat(configFile)
	if err != nil {
		return 0, errs.WrapUser(err, "unable to find file")
	}

	//Unmarshals based on file extension
	switch filepath.Ext(info.Name()) {
	case ".yml", ".yaml":
		err = yaml.Unmarshal(b, v)
	case ".json":
		err = json.Unmarshal(b, v)
	default:
		return 0, errs.NewUser("File type is not supported")
	}

	if err != nil {
		return 0, err
	}
	if v.Version == 0 {
		return 1, nil
	}
	return v.Version, nil
}

// UpgradeConfigVersion will convert a v1.Config to a v2.Config. Note that due to some semantic changes, this a lossy
//  conversion
func UpgradeConfigVersion(c1 *v1.Config) (*v2.Config, error) {
	c2 := &v2.Config{
		Version:  2,
		Docker:   c1.Docker,
		Accounts: map[string]v2.Account{},
		Envs:     map[string]v2.Env{},
	}

	def1 := c1.Defaults
	c2.Defaults = v2.Defaults{
		Common: v2.Common{
			Backend: &v2.Backend{
				Bucket:      util.StrPtr(def1.InfraBucket),
				DynamoTable: util.StrPtr(def1.InfraDynamoTable),
				Profile:     util.StrPtr(def1.AWSProfileBackend),
				Region:      util.StrPtr(def1.AWSRegionBackend),
			},
			ExtraVars: def1.ExtraVars,
			Providers: &v2.Providers{
				AWS: &v2.AWSProvider{
					AccountID:         &def1.AccountID,
					AdditionalRegions: def1.AWSRegions,
					Profile:           &def1.AWSProfileProvider,
					Region:            &def1.AWSRegionProvider,
					Version:           &def1.AWSProviderVersion,
				},
			},
			Owner:            util.StrPtr(def1.Owner),
			Project:          util.StrPtr(def1.Project),
			TerraformVersion: util.StrPtr(def1.TerraformVersion),
			Tools: &v2.Tools{
				TfLint:   def1.TfLint,
				TravisCI: c1.TravisCI,
			},
		},
	}

	for acctName, acct := range c1.Accounts {
		common := v2.Common{
			ExtraVars: acct.ExtraVars,
			Providers: &v2.Providers{
				AWS: &v2.AWSProvider{
					AccountID:         acct.AccountID,
					AdditionalRegions: acct.AWSRegions,
					Profile:           acct.AWSProfileProvider,
					Region:            acct.AWSRegionProvider,
					Version:           acct.AWSProviderVersion,
				},
			},
			Owner:            acct.Owner,
			Project:          acct.Project,
			TerraformVersion: acct.TerraformVersion,
		}
		if acct.InfraBucket != nil || acct.InfraDynamoTable != nil || acct.AWSProfileBackend != nil || acct.AWSRegionBackend != nil {
			common.Backend = &v2.Backend{
				Bucket:      acct.InfraBucket,
				DynamoTable: acct.InfraDynamoTable,
				Profile:     acct.AWSProfileBackend,
				Region:      acct.AWSRegionBackend,
			}
		}
		c2.Accounts[acctName] = v2.Account{
			Common: common,
		}
	}

	for envName, env := range c1.Envs {
		env2 := v2.Env{
			Common: v2.Common{
				ExtraVars: env.ExtraVars,
				Providers: &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID:         env.AccountID,
						AdditionalRegions: env.AWSRegions,
						Profile:           env.AWSProfileProvider,
						Region:            env.AWSRegionProvider,
						Version:           env.AWSProviderVersion,
					},
				},
				Owner:            env.Owner,
				Project:          env.Project,
				TerraformVersion: env.TerraformVersion,
			},
		}

		if env.InfraBucket != nil || env.InfraDynamoTable != nil || env.AWSProfileBackend != nil || env.AWSRegionBackend != nil {
			env2.Common.Backend = &v2.Backend{
				Bucket:      env.InfraBucket,
				DynamoTable: env.InfraDynamoTable,
				Profile:     env.AWSProfileBackend,
				Region:      env.AWSRegionBackend,
			}
		}

		env2.Components = map[string]v2.Component{}

		for componentName, component := range env.Components {
			c2 := v2.Component{
				Common: v2.Common{
					ExtraVars:        component.ExtraVars,
					Owner:            component.Owner,
					Project:          component.Project,
					TerraformVersion: component.TerraformVersion,
				},
				EKS:          component.EKS,
				Kind:         component.Kind,
				ModuleSource: component.ModuleSource,
			}

			if component.InfraBucket != nil || component.InfraDynamoTable != nil || component.AWSProfileBackend != nil || component.AWSRegionBackend != nil {
				c2.Backend = &v2.Backend{
					Bucket:      component.InfraBucket,
					DynamoTable: component.InfraDynamoTable,
					Profile:     component.AWSProfileBackend,
					Region:      component.AWSRegionBackend,
				}
			}

			if component.AccountID != nil || component.AWSRegions != nil || component.AWSProfileProvider != nil || component.AWSRegionProvider != nil || component.TerraformVersion != nil {
				c2.Providers = &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID:         component.AccountID,
						AdditionalRegions: component.AWSRegions,
						Profile:           component.AWSProfileProvider,
						Region:            component.AWSRegionProvider,
						Version:           component.AWSProviderVersion,
					},
				}
			}
			env2.Components[componentName] = c2
		}
		c2.Envs[envName] = env2
	}

	// we don't bother doing a deep copy bc we assume these structures are immutable
	c2.Modules = c1.Modules
	c2.Plugins = c1.Plugins

	return c2, nil
}
