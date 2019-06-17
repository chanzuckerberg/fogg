package yaml_migrate

import (
	"encoding/json"
	"io/ioutil"
	"os"

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

	logrus.Debug("Successfully Opened fogg.json")
	defer jsonFile.Close()

	//Read the fogg.json file, convert it to yml format
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return errs.WrapUser(err, "unable to read config file")
	}

	ymlData, err := jsonByteToYaml(byteValue)
	if err != nil {
		return errs.WrapUser(err, "unable to convert json to yaml")
	}

	return generateYMLFile(ymlData)
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

// Create YML file
// If file does not exist one will be made, otherwise the current
// yml file will be updated
func generateYMLFile(ymlData []byte) error {

	//Write creates a new file if one does not exist
	err := ioutil.WriteFile("fogg.yml", ymlData, 0644)

	return err
}

func OpenGitOrExit(fs afero.Fs) {
	_, err := fs.Stat(".git")
	if err != nil {
		// assuming this means no repository
		logrus.Fatal("fogg must be run from the root of a git repo")
		os.Exit(1)
	}
}
