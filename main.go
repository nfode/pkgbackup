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
	app             = kingpin.New("pkgbackup", "generic help")
	configFile      = app.Flag("configFile", "path to the configFile file").File()
	hostPackageFile = app.Flag("hostPackageFile", "alternative host packages").File()
	dryRun          = app.Flag("dryRun", "run the tool without installing anything").Bool()
	syncCommand     = app.Command("sync", "export package list")
	hostname        = os.Getenv("HOST")
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
	userInput := UserInput{dryRun: *dryRun, versionedPackages: versionedPackages, systemPackages: systemPackages, ignoredPackages: ignoredPackages}
	userInput.HandleHostPackagesChange()
	userInput.HandleSubscribedPackageChanges(hostConfig, config, baseDir)
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
	var packages []string
	if hostPackageFile != nil {
		packages, _ = ReadPackagesFromFile((*hostPackageFile).Name(), "")
	} else {
		cmd := exec.Command("yaourt", "-Qqe")
		result, err := cmd.Output()
		if err != nil {
			return []string{}, nil
		}
		resultAsString := string(result)
		packages = strings.Split(resultAsString, "\n")
	}
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
