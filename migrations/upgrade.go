package migrations

import (
	"github.com/chanzuckerberg/fogg/config"
	v1 "github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

//VersionUpgradeMigration Defines a fogg versioning upgrade
type VersionUpgradeMigration struct {
}

//Guard Checks the version of the config file and determines whether an upgrade is necessary
func (m *VersionUpgradeMigration) Guard(fs afero.Fs, configFile string) (bool, error) {
	_, version, err := config.FindConfig(fs, configFile)
	if err != nil {
		return false, err
	}

	switch version {
	case 1:
		return true, nil
	case 2:
		return false, nil
	default:
		return false, errs.NewUser("Config file version was not recognized")
	}
}

//Migrate Upgrades config file to most recent version
func (m *VersionUpgradeMigration) Migrate(fs afero.Fs, configFile string) (string, error) {
	bytes, version, err := config.FindConfig(fs, configFile)
	if err != nil {
		return "", err
	}
	switch version {
	case 1:
		c1, err := v1.ReadConfig(bytes)
		if err != nil {
			return "", err
		}
		c2, err := config.UpgradeConfigVersion(c1)
		if err != nil {
			return "", err
		}

		marshalled, err := yaml.Marshal(c2)
		if err != nil {
			return "", errs.WrapInternal(err, "Could not serialize config to yaml.")
		}
		err = afero.WriteFile(fs, configFile, marshalled, 0644)
		return configFile, errs.WrapInternal(err, "Could not write config to disk")

	case 2:
		logrus.Infof("config already v%d, nothing to do", version)
		return configFile, nil

	default:
		return configFile, errs.NewUserf("config version %d unrecognized", version)
	}
}
