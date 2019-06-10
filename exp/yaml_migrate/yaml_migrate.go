package yaml_migrate

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//ConvertToYaml method converts fogg.json to fogg.yml
func ConvertToYaml() error {
	jsonFile, err := os.Open("fogg.json")
	if err != nil {
		return err
	}
	logrus.Println("Successfully Opened fogg.json")
	defer jsonFile.Close()

	//Read the fogg.json file, convert it to yml format
	byteValue, readErr := ioutil.ReadAll(jsonFile)
	ymlData, yamlErr := jsonToYaml(byteValue)

	if readErr != nil {
		return readErr
	} else if yamlErr != nil {
		return yamlErr
	}

	generateYMLFile(ymlData)
	logrus.Println("Successfully created fogg.yml")
	return nil
}

//Convert an existing json file into yml text
func jsonToYaml(jsonFileData []byte) ([]byte, error) {
	var jsonObj interface{}

	//Convert jsonFileData into a generic interface representing json
	err := yaml.Unmarshal(jsonFileData, &jsonObj)
	if err != nil {
		return nil, err
	}

	// Convert generic json object into yml data
	return yaml.Marshal(jsonObj)
}

// Create YML file
// If file does not exist one will be made, otherwise the current
// yml file will be updated
//TODO: See if write overwrites a yaml file or adds to yaml file
func generateYMLFile(ymlData []byte) error {

	//Write creates a new file if one does not exist
	writeStatus := ioutil.WriteFile("fogg.yml", ymlData, 0644)
	return writeStatus
}
