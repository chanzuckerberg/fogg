package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type defaults struct {
	AWSRegion        string `json:"aws_region"`
	AWSProfile       string `json:"aws_profile"`
	InfraBucket      string `json:"infra_bucket"`
	Project          string `json:"project"`
	SharedInfraPath  string `json:"shared_infra_base"`
	TerraformVersion string `json:"terraform_version"`
}

type Config struct {
	Defaults defaults `json:"defaults"`
}

func ReadConfig(f io.ReadCloser) (*Config, error) {
	c := &Config{}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err2 := json.Unmarshal(b, c)
	if err2 != nil {
		return nil, err2
	}
	return c, nil
}
