package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

func main() {
	configFilePath := flag.String("configFile", "", "path to the config file")

	hostPackageFile := flag.String("hostPackageFile", "", "alternative host packages")
	dryRunFlag := flag.Bool("dryRun", false, "run the tool without changing anything")
	hostnameFlag := flag.String("hostname", "", "set custom hostname")
	flag.Parse()

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
