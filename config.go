package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Host struct {
	Name []string `yaml:"name"`
	File string   `yaml:"file"`
}
type Config struct {
	Hosts []Host `yaml:"host"`
}

func ParseConfigFile(file *os.File) Config {
	var config Config
	b, err := ioutil.ReadFile(file.Name())
	if err != nil {
		log.Fatal("error reading file" + err.Error())
		return Config{}
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("error unmarshaling yaml" + err.Error())
		return Config{}
	}

	return config
}
