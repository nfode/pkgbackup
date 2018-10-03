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
	app         = kingpin.New("pkgbackup", "generic help")
	configFile  = app.Flag("configFile", "path to the configFile file").File()
	syncCommand = app.Command("sync", "export package list")
	hostname    = os.Getenv("HOST")
)

func getHostConfig(hostname string, config Config) (Host, error) {
	for _, entry := range config.Hosts {
		if entry.Name == hostname {
			return entry, nil
		}
	}
	return Host{}, errors.New("file for hostname not found")
}

func ReadPackagesFromFile(fileName string, baseDir string) ([]string, error) {
	filePath := baseDir + "/" + fileName
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
	hostConfig, err := getHostConfig(hostname, config)
	versionedPackages, err := ReadPackagesFromFile(hostConfig.File, baseDir)
	systemPackages, err := GetInstalledPackages()
	if err != nil {
		fmt.Println("reading existing packages failed: " + err.Error())
	}
	ignoredPackages, err := ReadPackagesFromFile(hostConfig.IgnoreFile, baseDir)
	if err != nil {
		fmt.Println(fmt.Printf("Failed to read ignore file for host %v: %v", hostname, err.Error()))
	}
	comparisionResult := ComparePackages(versionedPackages, systemPackages)
	fmt.Println("Packages added to versioned packages:")
	for _, pkg := range comparisionResult.Added {
		fmt.Println(pkg)
	}
	for _, pkg := range comparisionResult.Removed {
		fmt.Println(pkg)
	}
	var outputs []OutputElement
	if hostConfig.SubscribeTo != nil {
		for _, subscribedTo := range hostConfig.SubscribeTo {
			subscribedToConfig, _ := getHostConfig(subscribedTo, config)
			subscribedToPackages, _ := ReadPackagesFromFile(subscribedToConfig.File, baseDir)
			comparisionResult := ComparePackages(subscribedToPackages, systemPackages)
			filteredResult := filterComparisonResult(ignoredPackages, comparisionResult)
			output := OutputElement{From: subscribedTo, ToInstall: filteredResult.Added, ToRemove: filteredResult.Removed}
			outputs = append(outputs, output)
		}
	}
	for _, output := range outputs {
		fmt.Println("Following packages were installed on: ", output.From)
		for _, pkg := range output.ToInstall {
			fmt.Println(pkg)
		}
		fmt.Println("Following packages were removed on: ", output.From)
		for _, pkg := range output.ToRemove {
			fmt.Println(pkg)
		}
	}
}
func filterComparisonResult(ignoredPackages []string, compareResult CompareResult) CompareResult {
	toRemove := clearList(compareResult.Removed, ignoredPackages)
	toAdd := clearList(compareResult.Added, ignoredPackages)
	return CompareResult{toAdd, toRemove, compareResult.Unchanged}
}

func clearList(toClear []string, ignoredPackages []string) []string {
	var result []string
	for _, value := range toClear {
		if contains(ignoredPackages, value) {
			continue
		} else {
			result = append(result, value)
		}
	}
	return result
}

func GetInstalledPackages() ([]string, error) {
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
