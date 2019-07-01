package migrations

import (
	"path/filepath"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

//JSONToYamlMigration Defines a fogg json to yaml config file migration
type JSONToYamlMigration struct {
}

//Guard Checks to see if config file needs to be converted to .yml
func (m *JSONToYamlMigration) Guard(fs afero.Fs, configFile string) (bool, error) {
	ext := filepath.Ext(configFile)

	switch ext {
	case ".json":
		return true, nil
	case ".yml", ".yaml":
		return false, nil
	default:
		return false, errs.NewUserf("Config file %s was not recognized", configFile)
	}
}

//Migrate method converts fogg.json to fogg.yml
func (m *JSONToYamlMigration) Migrate(fs afero.Fs, configFile string) (string, error) {
	logrus.Infof("Migrating %s to .yml file type", configFile)
	config, err := config.FindAndReadConfig(fs, configFile)
	if err != nil {
		return "", err
	}

	yml, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	err = afero.WriteFile(fs, "fogg.yml", yml, 0644)
	if err != nil {
		return "", err
	}
	logrus.Info("Created fogg.yml config file")

	err = fs.Remove(configFile)
	if err != nil {
		return "", err
	}
	logrus.Infof("Removed %s config file", configFile)
	return "fogg.yml", nil
}
