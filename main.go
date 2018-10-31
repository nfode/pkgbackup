package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func main() {

	configFilePath := flag.String("configFile", "", "path to the config file")
	hostPackageFile := flag.String("hostPackageFile", "", "alternative host packages")
	dryRunFlag := flag.Bool("dryRun", false, "run the tool without changing anything")
	hostnameFlag := flag.String("hostname", "", "set custom hostname")
	flag.CommandLine.Parse(os.Args[2:])

	command := os.Args[1]
	if command != "init" && command != "sync" {
		fmt.Printf("Command not specified.")
		os.Exit(1)
	}
	hostname, err := determineHostName(hostnameFlag)
	if err != nil {
		fmt.Println("Failed to determine hostname:", err)
		os.Exit(1)
	}
	if *configFilePath == "" {
		fmt.Println("Config file path not set!")
		os.Exit(1)
	}
	baseDir := path.Dir(*configFilePath)
	switch command {
	case "init":
		initConfiguration(hostname, configFilePath, baseDir)
		os.Exit(0)
	case "sync":
	}

	config := ParseConfigFile(*configFilePath)
	hostConfig, err := GetHostConfig(hostname, config)
	if err != nil {
		fmt.Printf("Host %v does not exist", hostname)
		os.Exit(1)
	}

	versionedPackages, err := ReadPackagesFromFile(hostConfig.File, baseDir)
	if err != nil {
		fmt.Println("reading versioned packages failed: " + err.Error())
	}

	systemPackages, err := GetInstalledPackages(hostPackageFile)
	if err != nil {
		fmt.Println("reading installed packages failed: " + err.Error())
	}

	ignoredPackages, err := ReadPackagesFromFile(hostConfig.IgnoreFile, baseDir)
	if err != nil {
		fmt.Println(fmt.Printf("Failed to read ignore file for host %v: %v", hostname, err.Error()))
	}

	userInput := UserInput{dryRun: *dryRunFlag, versionedPackages: versionedPackages, systemPackages: systemPackages, ignoredPackages: ignoredPackages}
	userInput.HandleHostPackagesChange()
	userInput.HandleSubscribedPackageChanges(hostConfig, config, baseDir)
	CommitPackageChanges(userInput, hostConfig, baseDir, *dryRunFlag)
}
func initConfiguration(hostname string, configFilePath *string, baseDir string) {
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		fmt.Printf("Creating directory went wrong: %v", err.Error())
	}

	ignoreFileName := hostname + "_ignore.txt"
	file := hostname + "_packages.txt"
	host := Host{Name: hostname, IgnoreFile: ignoreFileName, File: file, SubscribeTo: []string{}}
	config := Config{[]Host{host}}
	res, err := yaml.Marshal(config)
	if err != nil {
		fmt.Printf("Marshalling new config failed: %v", err.Error())
	}
	if err = ioutil.WriteFile(*configFilePath, res, CustomFileMode); err != nil {
		fmt.Printf("Writing new config failed: %v", err.Error())
	}

	packageFilePath := baseDir + "/" + file
	pkgs, err := GetInstalledPackages(nil)
	if err != nil {
		fmt.Printf("Failed to get current installed packages: %v", err.Error())
	}
	if err = ioutil.WriteFile(packageFilePath, []byte(strings.Join(pkgs, "\n")), CustomFileMode); err != nil {
		fmt.Printf("Failed to write installed packages: %v", err.Error())
	}
}

func determineHostName(hostnameFlag *string) (string, error) {
	hostname := *hostnameFlag
	if hostname == "" {
		host, err := os.Hostname()
		if err != nil {
			return "", err
		}
		hostname = host
	}
	return hostname, nil
}
