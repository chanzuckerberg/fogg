package util

import (
	"encoding/json"
	"io/ioutil"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

//ConvertToYaml method converts fogg.json to fogg.yml
func ConvertToYaml(fs afero.Fs, configFile string) error {
	jsonFile, err := fs.Open(configFile)
	if err != nil {
		return errs.WrapUser(err, "unable to open config file")
	}
	defer jsonFile.Close()

	logrus.Debug("Successfully Opened fogg.json")

	//Read the fogg.json file, convert it to yml format
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return errs.WrapUser(err, "unable to read config file")
	}

	ymlData, err := jsonByteToYaml(byteValue)
	if err != nil {
		return errs.WrapUser(err, "unable to convert json to yaml")
	}
	return ioutil.WriteFile("fogg.yml", ymlData, 0644)
}

//Convert an existing json file into yml text
func jsonByteToYaml(jsonFileData []byte) ([]byte, error) {
	var jsonObj interface{}

	//Creates a generic struct
	err := json.Unmarshal(jsonFileData, &jsonObj)
	if err != nil {
		return nil, err
	}

	// Converts generic struct to yaml output
	return yaml.Marshal(jsonObj)
}
