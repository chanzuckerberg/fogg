package migrations

import (
	"encoding/json"

	"github.com/chanzuckerberg/fogg/config"
	v1 "github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
	prompt "github.com/segmentio/go-prompt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

//VersionUpgradeMigration Defines a fogg version upgrade
type VersionUpgradeMigration struct {
}

//Description Describes the upgrade taking place
func (m *VersionUpgradeMigration) Description() string {
	return "v1 to v2 upgrade"
}

//Guard Checks the version of the config file and determines whether an upgrade is necessary
func (m *VersionUpgradeMigration) Guard(fs afero.Fs, configFile string) (bool, error) {
	_, version, err := config.FindConfig(fs, configFile)
	if err != nil {
		return false, err
	}

	switch version {
	case 1: //Upgrade is needed
		return true, nil
	case 2: //Latest version
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
	logrus.Infof("Attempting to upgrade v%d to the latest version", version)

	switch version {
	case 1: // Upgrades to v2 and craetes yaml file
		c1, err := v1.ReadConfig(bytes)
		if err != nil {
			return "", err
		}
		c2, err := config.UpgradeConfigVersion(c1)
		if err != nil {
			return "", err
		}

		marshalled, err := json.Marshal(c2)
		if err != nil {
			return "", errs.WrapInternal(err, "Could not serialize config to json.")
		}
		err = afero.WriteFile(fs, configFile, marshalled, 0644)
		return configFile, errs.WrapInternal(err, "Could not write config to disk")

	case 2: // Already v2, do nothing
		logrus.Infof("config already v%d, nothing to do", version)
		return configFile, nil

	default:
		return configFile, errs.NewUserf("config version %d was not recognized", version)
	}
}

//Prompt Checks to see if the user wants their version to be upgraded
func (m *VersionUpgradeMigration) Prompt() bool {
	return prompt.Confirm("Would you like to upgrade from v1 to v2?")
}
