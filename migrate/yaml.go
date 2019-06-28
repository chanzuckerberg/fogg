package migrate

import (
	"github.com/chanzuckerberg/fogg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

//ConvertToYaml method converts fogg.json to fogg.yml
func ConvertToYaml(fs afero.Fs, configFile string) error {
	config, err := config.FindAndReadConfig(fs, configFile)
	if err != nil {
		return err
	}

	yml, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = afero.WriteFile(fs, "fogg.yml", yml, 0644)
	if err != nil {
		return err
	}
	logrus.Info("Created fogg.yml config file")

	err = fs.Remove(configFile)
	if err != nil {
		return err
	}
	logrus.Infof("Removed %s config file", configFile)

	return nil
}
