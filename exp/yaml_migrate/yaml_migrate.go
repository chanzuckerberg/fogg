package yaml_migrate

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//TODO: Better Error Handling
//Converts fogg.json to fogg.yml
func JSONtoYML() error {
	jsonFile, err := os.Open("fogg.json")
	check(err)
	fmt.Println("Successfully Opened fogg.json")
	defer jsonFile.Close()

	//Read the fogg.json file, convert it to yml format
	byteValue, _ := ioutil.ReadAll(jsonFile)
	ymlData, _ := JSONtoYAML(byteValue)

	generateYMLFile(ymlData)
	fmt.Println("Successfully created fogg.yml")
	return nil
}

//Convert an existing json file into yml data
func JSONtoYAML(jsonFileData []byte) ([]byte, error) {
	var jsonObj interface{}

	//Convert jsonFileData into a generic interface representing json
	err := yaml.Unmarshal(jsonFileData, &jsonObj)
	if err != nil {
		return nil, err
	}

	return yaml.Marshal(jsonObj)
}

//TODO: Consider overwritteing yml file
// Create YML file
func generateYMLFile(ymlData []byte) {
	file, createErr := os.Create("fogg.yml")
	check(createErr)
	defer file.Close()

	writeErr := ioutil.WriteFile("fogg.yml", ymlData, 0644)
	check(writeErr)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
