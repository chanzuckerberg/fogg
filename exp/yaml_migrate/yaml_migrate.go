package yaml_migrate

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//TODO: Error Handle
//Converts fogg.json to fogg.yml
func JSONtoYML() error {
	dir, _ := os.Getwd()
	fmt.Printf("pwd = %v\n", dir)
	/* Open the file where data is being read and then close the file at the end of the function */
	jsonFile, err := os.Open("fogg.json")
	check(err)

	fmt.Println("Successfully Opened fogg.json")
	defer jsonFile.Close()

	//Read the fogg.json file, convert it to yml format

	byteValue, _ := ioutil.ReadAll(jsonFile)
	ymlData, _ := JSONtoYAML(byteValue)
	//Generate fogg.yml file
	generateYMLFile(ymlData)
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

	// Convert generic json object into yml data
	return yaml.Marshal(jsonObj)
}

// Create YML file
//TODO: Consider overwritteing yml file
func generateYMLFile(ymlData []byte) {
	dir, _ := os.Getwd()
	fmt.Printf("pwd = %v\n", dir)
	file, createErr := os.Create("fogg.yml")
	check(createErr)
	defer file.Close()

	fmt.Printf("File Created = %v", file)

	writeErr := ioutil.WriteFile("fogg.yml", ymlData, 0644)
	check(writeErr)
}

//Throws exception TODO:
func check(e error) {
	if e != nil {
		panic(e)
	}
}
