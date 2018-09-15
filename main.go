package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	hostname, _ = os.Hostname()
	app         = kingpin.New("pkgbackup", "generic help")
	configFile  = app.Flag("configFile", "path to the configFile file").File()
	syncCommand = app.Command("sync", "export package list")
)

func getFileForHostname(hostname string, config Config) (string, error) {
	for _, entry := range config.Hosts {
		for _, name := range entry.Name {
			if name == hostname {
				return entry.File, nil
			}
		}
	}
	return string(""), errors.New("file for hostname not found")
}

func ReadPackagesFromFile(config Config, hostname string, baseDir string) ([]string, error) {
	filePath, err := getFileForHostname(hostname, config)
	filePath = baseDir + "/" + filePath
	if err != nil {
		return []string{}, err
	}
	data, fileErr := ioutil.ReadFile(filePath)
	if fileErr != nil {
		return []string{}, fileErr
	}

	text := string(data[:])

	packages := strings.Split(text, "\n")

	return packages, nil
}

func main() {
	res := kingpin.MustParse(app.Parse(os.Args[1:]))
	config := ParseConfigFile(*configFile)
	switch res {
	case syncCommand.FullCommand():
		sync(config, path.Dir((*configFile).Name()))
	}
}
func sync(config Config, baseDir string) {
	exportedPackages, err := ReadPackagesFromFile(config, hostname, baseDir)
	systemPackages, err := GetSystemPackages()
	if err != nil {
		fmt.Println("reading existing packages failed: " + err.Error())
	}
	comparisionResult := ComparePackages(exportedPackages, systemPackages)

	AskUser(comparisionResult)
}

func GetSystemPackages() ([]string, error) {
	cmd := exec.Command("yaourt", "-Qqe")
	result, err := cmd.Output()
	if err != nil {
		return []string{}, nil
	}
	resultAsString := string(result)
	packages := strings.Split(resultAsString, "\n")
	return clearPackages(packages), nil
}

func clearPackages(packages []string) []string {
	var clearedPackages []string
	for _, v := range packages {
		if len(v) > 1 {
			clearedPackages = append(clearedPackages, v)
		}
	}
	return clearedPackages
}
