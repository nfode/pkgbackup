package main

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Host struct {
	Name        string   `yaml:"name"`
	File        string   `yaml:"file"`
	IgnoreFile  string   `yaml:"ignore"`
	SubscribeTo []string `yaml:"subscribeTo,omitempty"`
}
type Config struct {
	Hosts []Host `yaml:"hosts"`
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

func GetHostConfig(hostname string, config Config) (Host, error) {
	for _, entry := range config.Hosts {
		if entry.Name == hostname {
			return entry, nil
		}
	}
	return Host{}, errors.New("file for hostname not found")
}

func ReadPackagesFromFile(fileName string, baseDir string) ([]string, error) {
	filePath := fileName
	if baseDir != "" {
		filePath = baseDir + "/" + fileName
	}
	data, fileErr := ioutil.ReadFile(filePath)
	if fileErr != nil {
		return []string{}, fileErr
	}

	text := string(data[:])

	packages := strings.Split(text, "\n")

	return packages, nil
}
