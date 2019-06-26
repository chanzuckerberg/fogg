package config

import (
	v1 "github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// Upgrade applies in-place upgrades to a configFile
func Upgrade(fs afero.Fs, configFile string) error {
	bytes, version, err := FindConfig(fs, configFile)
	if err != nil {
		return err
	}
	switch version {
	case 1:
		c1, err := v1.ReadConfig(bytes)
		if err != nil {
			return err
		}
		c2, err := UpgradeConfigVersion(c1)
		if err != nil {
			return err
		}

		marshalled, err := yaml.Marshal(c2)
		if err != nil {
			return errs.WrapInternal(err, "Could not serialize config to yaml.")
		}
		err = afero.WriteFile(fs, configFile, marshalled, 0644)
		return errs.WrapInternal(err, "Could not write config to disk")
	case 2:
		logrus.Infof("config already v%d, nothing to do", version)
		return nil

	default:
		return errs.NewUserf("config version %d unrecognized", version)
	}
}
